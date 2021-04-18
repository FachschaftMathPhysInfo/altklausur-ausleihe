<template>
  <div>
    <v-app-bar color="primary" dense dark>
      <v-app-bar-nav-icon></v-app-bar-nav-icon>

      <v-toolbar-title>Altklausurausleihe</v-toolbar-title>

      <v-spacer></v-spacer>
      <v-spacer></v-spacer>

      <v-text-field
        v-model="search"
        prepend-inner-icon="mdi-magnify"
        label="Filter Altklausuren, z.B. nach Prüfenden oder Veranstaltungen"
        single-line
        hide-details
        clearable
      ></v-text-field>
    </v-app-bar>
    <v-container>
      <v-row v-if="search">
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
        v-model="selected"
        :headers="headers"
        :items="exams"
        :items-per-page="-1"
        :search="search"
        :hide-default-footer="true"
        show-expand
        show-select
      >
        <template v-slot:[`item.subject`]="{ item }">
          <v-chip v-if="item.subject"
            ><v-icon :color="getSubjectColor(item.subject)" left>{{
              getSubjectIcon(item.subject)
            }}</v-icon
            >{{ item.subject }}</v-chip
          >
        </template>
        <template v-slot:expanded-item="{ headers, item }">
          <td :colspan="headers.length">
            <iframe :src="item.path" style="width:100%; height: 1500px;" />
          </td>
        </template>
        <template v-slot:no-data>
          <span
            >Es wurden keine Klausuren passend zu den Suchkriterien
            gefunden!</span
          >
        </template>
      </v-data-table>
      <v-tooltip left v-if="selected.length > 0">
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
            @click="downloadExams"
            ><v-icon>mdi-download</v-icon></v-btn
          >
        </template>
        <span>Ausgewählte Altklausuren herunterladen</span>
      </v-tooltip>
      <v-tooltip left v-if="selected.length == 0">
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
        <span>Wähle eine Altklausur aus, um sie dann herunterzuladen.</span>
      </v-tooltip>
    </v-container>
  </div>
</template>

<script>
import testexams from "@/assets/testexams.json";
export default {
  name: "ExamList",
  components: {},
  data() {
    const self = this;
    return {
      selected: [],
      search: "",
      examiner: null,
      moduleName: null,
      subjects: ["Mathe", "Physik", "Info"],
      fromSemester: null,
      toSemester: null,
      exams: testexams,
      headers: [
        { text: "", value: "data-table-expand" },
        {
          text: "Veranstaltung",
          value: "moduleName"
        },
        { text: "Prüfer", value: "examiners" },
        {
          text: "Semester",
          value: "semester",
          sortable: true,
          sort: (a, b) => self.semesterBefore(a, b)
        },
        { text: "Fach", value: "subject" },
        { text: "", value: "data-table-select" }
      ]
    };
  },
  computed: {
    semesters() {
      return testexams
        .map(exam => ({ name: exam.semester }))
        .sort((a, b) => this.semesterBefore(a.name, b.name));
    }
  },

  methods: {
    downloadExams() {
      alert("To be implemented: download selected PDFs.");
    },
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
      this.exams = testexams.filter(
        exam =>
          this.subjects.includes(exam.subject) &&
          (this.moduleName == null ||
            exam.moduleName.includes(this.moduleName)) &&
          (this.examiner == null || exam.examiners.includes(this.examiner)) &&
          (this.fromSemester == null ||
            this.semesterBeforeOrEqual(this.fromSemester, exam.semester)) &&
          (this.toSemester == null ||
            this.semesterBeforeOrEqual(exam.semester, this.toSemester))
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
    }
  }
};
</script>