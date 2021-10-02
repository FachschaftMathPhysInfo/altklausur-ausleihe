<template>
  <div>
    <v-container>
      <v-row v-if="this.$parent.search">
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
      </v-row>
      <v-data-table
        :headers="headers"
        :items="exams"
        item-key="UUID"
        :items-per-page="-1"
        :search="this.$parent.search"
        :hide-default-footer="true"
        show-expand
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
            <v-icon>
              mdi-download
            </v-icon>
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
              style="width: 100%; height: 1500px;"
            />
          </td>
        </template>
        <template v-slot:no-data>
          <span>{{ $t("examlist.no_exams_found") }}</span>
        </template>
      </v-data-table>
      <v-tooltip left>
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            v-bind="attrs"
            v-on="on"
            elevation="2"
            fixed
            right
            bottom
            color="primary"
            fab
            @click="help()"
            ><v-icon>mdi-help</v-icon></v-btn
          >
        </template>
        <span>{{ $t("examlist.click_for_explanation") }}</span>
      </v-tooltip>
    </v-container>
    {{ examiner }}
  </div>
</template>

<script>
import gql from "graphql-tag";

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

export default {
  name: "ExamList",
  components: {},
  data() {
    // const self = this;
    return {
      examiner: null,
      moduleName: null,
      subjects: ["Mathe", "Physik", "Info"],
      fromSemester: null,
      toSemester: null,
      exams: [],
      originalExams: [],
    };
  },
  computed: {
    headers() {
      return [
        { text: "", value: "data-table-expand" },
        {
          text: this.$t("examlist.module"),
          value: "moduleName",
        },
        { text: this.$t("examlist.examiner"), value: "examiners" },
        {
          text: "Semester",
          value: "combinedSemester",
          sortable: true,
          sort: (a, b) => this.semesterBefore(a, b),
        },
        { text: this.$t("examlist.subject"), value: "subject" },
        { text: this.$t("examlist.download"), value: "download" },
      ];
    },
    semesters() {
      if (this.exams.length > 0) {
        return this.exams
          .filter((exam) => exam.combinedSemester.trim() != "")
          .map((exam) => ({ name: exam.combinedSemester }))
          .sort((a, b) => this.semesterBefore(a.name, b.name));
      } else {
        return [];
      }
    },
  },

  methods: {
    help() {
      alert(
        "To be implemented: Open help dialog with very detailed instructions"
      );
    },
    disableFromSemester(semester) {
      if (this.toSemester == null) return false;
      return this.semesterBefore(this.toSemester, semester.name);
    },
    disableToSemester(semester) {
      if (this.fromSemester == null) return false;
      return this.semesterBefore(semester.name, this.fromSemester);
    },
    semesterBefore(thisSemester, otherSemester) {
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
    },
    semesterBeforeOrEqual(thisSemester, otherSemester) {
      return (
        this.semesterBefore(thisSemester, otherSemester) ||
        thisSemester == otherSemester
      );
    },
    filterExams() {
      if (this.originalExams.length == 0) {
        this.originalExams = this.exams;
      }
      this.exams = this.originalExams.filter(
        (exam) =>
          this.subjects.includes(exam.subject) &&
          (this.moduleName == null ||
            exam.moduleName.includes(this.moduleName)) &&
          (this.examiner == null || exam.examiners.includes(this.examiner)) &&
          (this.fromSemester == null ||
            this.semesterBeforeOrEqual(
              this.fromSemester,
              exam.combinedSemester
            )) &&
          (this.toSemester == null ||
            this.semesterBeforeOrEqual(exam.combinedSemester, this.toSemester))
      );
    },
    getSubjectColor(subject) {
      if (subject == "Mathe") {
        return "green";
      } else if (subject == "Physik") {
        return "blue";
      } else if (subject == "Info") {
        return "orange";
      } else {
        return "gray";
      }
    },
    getSubjectIcon(subject) {
      if (subject == "Mathe") {
        return "mdi-android-studio";
      } else if (subject == "Physik") {
        return "mdi-atom";
      } else if (subject == "Info") {
        return "mdi-laptop";
      } else {
        return "mdi-label";
      }
    },
    async downloadAltklausur(exam) {
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
    },
    async getMarkedExamURLFromRow(row) {
      // retrieve urls from backend when exam row is opened
      if (!row.item.viewUrl) {
        await this.watermarkExam(row.item.UUID);
        await this.getExamURLs(row.item, false);
      }
    },
    async watermarkExam(UUID) {
      // Call to the graphql mutation to initiate watermarking process in backend
      await new Promise((f) => setTimeout(f, 500));
      await this.$apollo.mutate({
        mutation: gql`
          mutation($UUID: String!) {
            requestMarkedExam(StringUUID: $UUID)
          }
        `,
        variables: {
          UUID: UUID,
        },
      });
    },
    async getExamURLs(exam, openDownload) {
      // Call to the graphql query, to retrieve URLs of exam PDFs. Repeat 5 times if not successful and then time out
      for (let i = 0; i < 5; i++) {
        const result = await this.$apollo.query({
          query: gql`
            query($UUID: String!) {
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
            this.openExam(exam.downloadUrl);
          }
          break;
        }
      }
      if (exam.loading) {
        // request failed even after 5 retries
        alert(this.$t("examlist.request_failed"));
      }
    },
    openExam(url) {
      const link = document.createElement("a");
      link.href = url;
      link.click();
    },
  },
  apollo: {
    exams: {
      query: EXAMS_QUERY,
      update: (data) => {
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
      },
    },
  },
};
</script>
