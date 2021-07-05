<template>
  <div>
    <v-container>
      <v-row v-if="this.$parent.search">
        <v-col sm="2">
          <v-text-field
            v-model="moduleName"
            prepend-inner-icon="mdi-book-open-variant"
            label="Veranstaltungsname"
            hint="Auch Abkürzungen und Varianten werden berücksichtigt"
            single-line
            clearable
            @input="filterExams"
          ></v-text-field>
        </v-col>
        <v-col sm="2">
          <v-text-field
            v-model="examiner"
            prepend-inner-icon="mdi-account"
            label="Prüfende eingrenzen"
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
            label="ab Semester"
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
            label="bis Semester"
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
              <span class="hidden-sm-and-down">Mathematik</span>
            </v-btn>

            <v-btn value="Physik">
              <v-icon :color="getSubjectColor('Physik')" left>
                mdi-atom
              </v-icon>
              <span class="hidden-sm-and-down">Physik</span>
            </v-btn>

            <v-btn value="Info">
              <v-icon :color="getSubjectColor('Info')" left>
                mdi-laptop
              </v-icon>
              <span class="hidden-sm-and-down">Informatik</span>
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
          <a v-if="item.downloadUrl" :href="item.downloadUrl" target="_blank">
            <tooltipped-icon
              icon="mdi-download"
              color="green"
              text="Altklausur herunterladen"
              position="bottom"
              @clicked="downloadExam(item.downloadUrl)"
            ></tooltipped-icon>
          </a>
          <tooltipped-icon
            v-if="!item.downloadUrl"
            icon="mdi-stamper"
            color="primary"
            text="Altklausur mit Wasserzeichen versehen"
            position="bottom"
            @clicked="getMarkedExamURL(item)"
          ></tooltipped-icon>
        </template>
        <template v-slot:expanded-item="{ headers, item }">
          <td :colspan="headers.length">
            <p v-if="!item.viewUrl">Watermarking and Loading Exam ...</p>
            <iframe
              v-if="item.viewUrl"
              :src="item.viewUrl"
              style="width: 100%; height: 1500px;"
            />
          </td>
        </template>
        <template v-slot:no-data>
          <span
            >Es wurden keine Klausuren passend zu den Suchkriterien
            gefunden!</span
          >
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
        <span>Klicke hier für eine Anleitung</span>
      </v-tooltip>
    </v-container>
  </div>
</template>

<script>
import TooltippedIcon from "./generic/TooltippedIcon.vue";
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
  components: { TooltippedIcon },
  data() {
    const self = this;
    return {
      examiner: null,
      moduleName: null,
      subjects: ["Mathe", "Physik", "Info"],
      fromSemester: null,
      toSemester: null,
      exams: [],
      headers: [
        { text: "", value: "data-table-expand" },
        {
          text: "Veranstaltung",
          value: "moduleName",
        },
        { text: "Prüfer", value: "examiners" },
        {
          text: "Semester",
          value: "combinedSemester",
          sortable: true,
          sort: (a, b) => self.semesterBefore(a, b),
        },
        { text: "Fach", value: "subject" },
        { text: "Download", value: "download" },
      ],
    };
  },
  computed: {
    semesters() {
      if (this.exams.length > 0) {
        console.log(this.exams);
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
      this.exams = this.exams.filter(
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
    async watermarkExam(UUID) {
      // Call to the graphql mutation
      const result = await this.$apollo.mutate({
        mutation: gql`
          mutation($UUID: String!) {
            requestMarkedExam(StringUUID: $UUID)
          }
        `,
        variables: {
          UUID: UUID,
        },
      });
      if (!result) {
        // this seems to be necessary to watermark new exams
        console.log(result);
      }
    },
    async getExamURLs(exam) {
      // Call to the graphql query
      const result = await this.$apollo.query({
        // Query
        query: gql`
          query($UUID: String!) {
            getExam(StringUUID: $UUID) {
              viewUrl
              downloadUrl
            }
          }
        `,
        // Parameters
        variables: {
          UUID: exam.UUID,
        },
      });
      exam.viewUrl = result.data.getExam.viewUrl;
      exam.downloadUrl = result.data.getExam.downloadUrl;

      this.$forceUpdate();
    },
    async getMarkedExamURLFromRow(row) {
      console.log(row.item);
      await this.getMarkedExamURL(row.item);
    },
    async getMarkedExamURL(exam) {
      console.log(exam);
      await this.watermarkExam(exam.UUID);
      await this.getExamURLs(exam);
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
        });

        return data.exams;
      },
    },
  },
};
</script>
