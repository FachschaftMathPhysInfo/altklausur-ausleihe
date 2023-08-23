<script setup>
import { ref, computed, provide } from "vue";
import { gql } from "@apollo/client/core";
import {
  provideApolloClient,
  useQuery,
  useResult,
} from "@vue/apollo-composable";
import { ApolloClient, InMemoryCache } from "@apollo/client/core";
import { useI18n } from "vue-i18n";

// fetch all exams
const EXAMS_QUERY = gql`
  query {
    exams {
      UUID
      subject
      moduleName
      moduleAltName
      year
      examiners
      semester
    }
  }
`;

// Name of the localStorage item
const AUTH_TOKEN = "apollo-token";

// Http endpoint
const httpEndpoint = "http://localhost:8081/query";
// const httpEndpoint = 'http://localhost:8081/query'

// ws endpoint
// const wsEndpoint = SERVER_WS || 'ws://localhost:8081/query'
// const wsEndpoint = 'ws://localhost:8081/query'

// Config
const defaultOptions = {
  // You can use `https` for secure connection (recommended in production)
  httpEndpoint,
  // You can use `wss` for secure connection (recommended in production)
  // Use `null` to disable subscriptions
  //   wsEndpoint,
  // LocalStorage token
  tokenName: AUTH_TOKEN,
  // Enable Automatic Query persisting with Apollo Engine
  persisting: false,
  // Use websockets for everything (no HTTP)
  // You need to pass a `wsEndpoint` for this to work
  websocketsOnly: false,
  // Is being rendered on the server?
  ssr: false,

  // Override default apollo link
  // note: don't override httpLink here, specify httpLink options in the
  // httpLinkOptions property of defaultOptions.

  // Override default cache
  cache: new InMemoryCache(),

  // Override the way the Authorization header is set
  // getAuth: (tokenName) => ...

  // Additional ApolloClient options
  // apollo: { ... }

  // Client local data (see apollo-link-state)
  // clientState: { resolvers: { ... }, defaults: { ... } }
};
const apolloClient = new ApolloClient(defaultOptions);

provideApolloClient(apolloClient);

const { result, loading, error } = useQuery(EXAMS_QUERY);

const repositories = useResult(result, [], (data) => {
  data.exams.forEach((exam) => {
    // Set undefined elements to empty strings
    Object.keys(exam).forEach((key) => {
      exam[key] = exam[key] ? exam[key] : "";
    });

    // combine year and semester to combined semester
    exam.combinedSemester = `${exam.semester} ${exam.year}`;
    exam.loading = null;
    exam.disabled = null;
  });

  return data.exams;
});

// UI
const notAuthenticatedDialog = ref(false);
const helpDialog = ref(false);

// filtering
const examiner = ref(null);
const moduleName = ref(null);
const subjects = ["Mathe", "Physik", "Info"];
const fromSemester = ref(null);
const toSemester = ref(null);
// const exams = ref([]);
const originalExams = ref([]);

const i18n = useI18n();

const headers = computed(() => {
  return [
    { title: "", key: "data-table-expand" },
    {
      title: i18n.t("examlist.module"),
      key: "moduleName",
    },
    { title: i18n.t("examlist.examiner"), key: "examiners" },
    {
      title: "Semester",
      key: "combinedSemester",
      sortable: true,
      sort: (a, b) => this.semesterSort(a, b),
    },
    { title: i18n.t("examlist.subject"), key: "subject" },

    { title: i18n.t("examlist.download"), key: "download" },
  ];
});

const semesters = computed(() => {
  if (exams.length > 0) {
    return exams
      .filter((exam) => exam.combinedSemester.trim() != "")
      .map((exam) => ({ name: exam.combinedSemester }))
      .sort((a, b) => this.semesterSort(a.name, b.name));
  } else {
    return [];
  }
});

function openNotAuthenticatedDialog() {
  this.notAuthenticatedDialog = true;
}

function help() {
  this.helpDialog = true;
}

function disableFromSemester(semester) {
  if (this.toSemester == null) return false;
  return this.semesterBefore(this.toSemester, semester.name);
}

function disableToSemester(semester) {
  if (this.fromSemester == null) return false;
  return this.semesterBefore(semester.name, this.fromSemester);
}
function semesterSort(thisSemester, otherSemester) {
  // splits semester labels into year (index 1) and season (index 0)
  const thisSem = thisSemester.split(" ");
  const otherSem = otherSemester.split(" ");
  if (thisSem[1] < otherSem[1]) {
    return 1;
  } else if (thisSem[1] == otherSem[1]) {
    if (thisSem[0] < otherSem[0]) {
      return 1;
    }
  }
  return -1;
}

function semesterBefore(thisSemester, otherSemester) {
  // splits semester labels into year (index 1) and season (index 0)
  const thisSem = thisSemester.split(" ");
  const otherSem = otherSemester.split(" ");
  if (thisSem[1] < otherSem[1]) {
    return true;
  } else if (thisSem[1] == otherSem[1]) {
    if (thisSem[0] < otherSem[0]) {
      return true;
    }
  }
  return false;
}

function semesterBeforeOrEqual(thisSemester, otherSemester) {
  return (
    this.semesterBefore(thisSemester, otherSemester) ||
    thisSemester == otherSemester
  );
}

function filterExams() {
  if (this.originalExams.length == 0) {
    this.originalExams = this.exams;
  }
  this.exams = this.originalExams.filter(
    (exam) =>
      this.subjects.includes(exam.subject) &&
      (this.moduleName == null || exam.moduleName.includes(this.moduleName)) &&
      (this.examiner == null || exam.examiners.includes(this.examiner)) &&
      (this.fromSemester == null ||
        this.semesterBeforeOrEqual(this.fromSemester, exam.combinedSemester)) &&
      (this.toSemester == null ||
        this.semesterBeforeOrEqual(exam.combinedSemester, this.toSemester))
  );
}

function getSubjectColor(subject) {
  if (subject == "Mathe") {
    return "green";
  } else if (subject == "Physik") {
    return "blue";
  } else if (subject == "Info") {
    return "orange";
  } else {
    return "gray";
  }
}

function getSubjectIcon(subject) {
  if (subject == "Mathe") {
    return "mdi-android-studio";
  } else if (subject == "Physik") {
    return "mdi-atom";
  } else if (subject == "Info") {
    return "mdi-laptop";
  } else {
    return "mdi-label";
  }
}

async function downloadAltklausur(exam) {
  // download the exam in a two step process: 1. watermark 2. get URLs
  if (!exam.downloadUrl) {
    exam.loading = true;
    exam.disabled = true;

    await this.watermarkExam(exam.UUID);
    await this.getExamURLs(exam, true);

    exam.loading = false;
    exam.disabled = false;
  } else {
    // simply open exam if it has been processed already
    this.openExam(exam.downloadUrl);
  }
}
async function getMarkedExamURLFromRow(row) {
  // retrieve urls from backend when exam row is opened
  if (!row.item.viewUrl) {
    await this.watermarkExam(row.item.UUID);
    await this.getExamURLs(row.item, false);
  }
}

async function watermarkExam(UUID) {
  // Call to the graphql mutation to initiate watermarking process in backend
  await new Promise((f) => setTimeout(f, 500));
  await this.$apollo.mutate({
    mutation: gql`
      mutation ($UUID: String!) {
        requestMarkedExam(StringUUID: $UUID)
      }
    `,
    variables: {
      UUID: UUID,
    },
  });
}

async function getExamURLs(exam, openDownload) {
  // Call to the graphql query, to retrieve URLs of exam PDFs. Repeat 5 times if not successful and then time out
  for (let i = 0; i < 5; i++) {
    const result = await this.$apollo.query({
      query: gql`
        query ($UUID: String!) {
          getExam(StringUUID: $UUID) {
            viewUrl
            downloadUrl
          }
        }
      `,
      variables: {
        UUID: exam.UUID,
      },
      fetchPolicy: "network-only", //necessary, otherwise result.data.getExam will always be null
    });

    if (result.data.getExam == null) {
      // watermarked result isn't ready yet => wait a moment and retry
      await new Promise((f) => setTimeout(f, 1000));
    } else {
      exam.viewUrl = result.data.getExam.viewUrl;
      exam.downloadUrl = result.data.getExam.downloadUrl;
      this.$forceUpdate();
      if (openDownload) {
        openExam(exam.downloadUrl);
      }
      break;
    }
  }
  if (exam.loading) {
    // request failed even after 5 retries
    alert(this.$t("examlist.request_failed"));
  }
}

function openExam(url) {
  const link = document.createElement("a");
  link.href = url;
  link.click();
}

function isMobile() {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
    navigator.userAgent
  );
}
</script>

<template>
  <div>
    <v-container>
      <!-- v-row v-if="this.$parent.search">
        <v-col sm="2">
          <v-text-field
            v-model="moduleName"
            prepend-inner-icon="mdi-book-open-variant"
            :label="$t('examlist.eventname')"
            :hint="$t('examlist.hint')"
            single-line
            clearable
            @input="filterExams"
            @change="filterExams"
          ></v-text-field>
        </v-col>
        <v-col sm="2">
          <v-text-field
            v-model="examiner"
            prepend-inner-icon="mdi-account"
            :label="$t('examlist.filter_lecturers')"
            single-line
            clearable
            @input="filterExams"
          ></v-text-field>
        </v-col>
        <v-col sm="2">
          <v-select
            v-model="fromSemester"
            :items="semesters"
            item-text="name"
            :item-disabled="disableFromSemester"
            :label="$t('examlist.from_semester')"
            clearable
            @change="filterExams"
          ></v-select>
        </v-col>
        <v-col sm="2">
          <v-select
            v-model="toSemester"
            :items="semesters"
            item-text="name"
            :item-disabled="disableToSemester"
            :label="$t('examlist.to_semester')"
            clearable
            @change="filterExams"
          ></v-select>
        </v-col>
        <v-col sm="4" align="center">
          <v-btn-toggle
            v-model="subjects"
            multiple
            rounded
            @change="filterExams"
          >
            <v-btn value="Mathe">
              <v-icon :color="getSubjectColor('Mathe')" left>
                mdi-android-studio
              </v-icon>
              <span class="hidden-sm-and-down">{{ $t("examlist.maths") }}</span>
            </v-btn>

            <v-btn value="Physik">
              <v-icon :color="getSubjectColor('Physik')" left>
                mdi-atom
              </v-icon>
              <span class="hidden-sm-and-down">{{
                $t("examlist.physics")
              }}</span>
            </v-btn>

            <v-btn value="Info">
              <v-icon :color="getSubjectColor('Info')" left>
                mdi-laptop
              </v-icon>
              <span class="hidden-sm-and-down">{{
                $t("examlist.computer_science")
              }}</span>
            </v-btn>
          </v-btn-toggle>
        </v-col>
      </v-row> -->
      <!-- :search="this.$parent.search" -->

      <!-- <v-data-table
        :items-per-page="-1"
        :headers="headers"
        :items="exams"
        item-value="name"
      ></v-data-table> -->

      <p v-if="loading">Loding...</p>
      <p>Error: {{ error }}</p>
      <v-data-table
        :headers="headers"
        :items="exams"
        item-key="UUID"
        :items-per-page="-1"
        :hide-default-footer="true"
        :show-expand="!isMobile()"
        @item-expanded="getMarkedExamURLFromRow"
      >
        <template v-slot:[`item.subject`]="{ item }">
          <v-chip v-if="item.subject"
            ><v-icon :color="getSubjectColor(item.subject)" left>{{
              getSubjectIcon(item.subject)
            }}</v-icon
            >{{ item.subject }}</v-chip
          >
        </template>
        <template v-slot:[`item.download`]="{ item }">
          <v-btn
            :loading="item.loading"
            :disabled="item.disabled"
            color="primary"
            @click="downloadAltklausur(item)"
            rounded
          >
            <v-icon> mdi-download </v-icon>
            {{ $t("examlist.downloaden") }}
          </v-btn>
        </template>
        <template v-slot:expanded-item="{ headers, item }">
          <td :colspan="headers.length">
            <div v-if="!item.viewUrl" class="text-center">
              <h4>{{ $t("examlist.watermarking") }}</h4>
              <v-progress-circular
                indeterminate
                color="primary"
              ></v-progress-circular>
            </div>

            <iframe
              v-if="item.viewUrl"
              :src="item.viewUrl"
              style="width: 100%; height: 1500px"
            />
          </td>
        </template>

        <template v-slot:no-data>
          <span>{{ $t("examlist.no_exams_found") }}</span>
        </template>
      </v-data-table>

      <v-layout-item model-value position="bottom" class="text-end" size="88">
        <div class="ma-4">
          <v-tooltip
            :text="$t('examlist.click_for_explanation')"
            location="start"
          >
            <template v-slot:activator="{ props }">
              <v-btn
                v-bind="props"
                icon="mdi-help"
                size="large"
                color="primary"
                elevation="2"
                @click="help()"
              />
            </template>
          </v-tooltip>
        </div>
      </v-layout-item>

      <v-dialog
        v-model="notAuthenticatedDialog"
        transition="dialog-bottom-transition"
        max-width="600"
      >
        <v-card>
          <v-toolbar color="primary" dark>
            <v-icon class="pr-3" large>mdi-alert</v-icon>
            {{ $t("auth.not_authenticated") }}
          </v-toolbar>
          <v-card-text>
            <div class="text pa-6">
              {{ $t("auth.text") }}
            </div>
          </v-card-text>
          <v-card-actions class="justify-end">
            <v-btn
              depressed
              @click="dialog.value = false"
              color="primary"
              elevation="2"
              href="https://moodle.uni-heidelberg.de/mod/lti/view.php?id=547942"
            >
              {{ $t("auth.login") }}
            </v-btn>
            <v-btn text @click="notAuthenticatedDialog.value = false">{{
              $t("auth.close")
            }}</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

      <v-dialog
        v-model="helpDialog"
        transition="dialog-bottom-transition"
        max-width="900"
      >
        <v-card>
          <v-toolbar color="primary" dark>
            <v-icon class="pr-3" large>mdi-help</v-icon>
            {{ $t("help.help") }}
          </v-toolbar>
          <v-card-text class="pa-6">
            <div
              v-for="(item, i) in $t('help.questions_and_answers')"
              v-bind:key="i"
            >
              <div class="text-h6 pa-2">{{ i + 1 }}. {{ item.question }}</div>
              <div class="text pa-2">
                {{ item.answer }}
              </div>
            </div>
          </v-card-text>
          <v-card-actions class="justify-end">
            <v-btn text @click="helpDialog.value = false">{{
              $t("help.close")
            }}</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-container>
  </div>
</template>

<style scoped></style>
