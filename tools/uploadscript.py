import json
from datetime import datetime
import os.path


def upload_report(JWT_TOKEN, report, filepath):

    # Abort if any argument is not given
    if not os.path.isfile(filepath):
        print(f"Error, file \"{filepath}\" not found!")
        return

    exec_cmd = f"./upload_one.sh {JWT_TOKEN} \'{json.dumps(report)}\' {filepath}"
    # print(exec_cmd)
    os.system(exec_cmd)


def norm_subject_name(input_name):
    for subject in ["Mathe", "Physik", "Info"]:
        if subject.lower() in input_name.lower():
            return subject


def date2semester(datestring):
    """
    Sample input: "2020-02-04T23:00:00.000Z"
    returns: (Semester, Year)
    """
    date = datetime.strptime(datestring, '%Y-%m-%dT%H:%M:%S.%fZ')
    # November to April
    if date.month <= 4 or 11 <= date.month:
        return ("WiSe", date.year)
    else:
        # May to October
        return ("SoSe", date.year)

    return (None, None)


def main():
    # Enter the JWT Token
    JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjEyMCIsImV4cCI6MTYzOTAwNDkwM30.bdGwQQB1XJ5lzpZiuvNxvLLhR85NzV1sysleA-1i0jo"

    # load the inputs
    with open("./reports_raw.json") as inputfile:
        input_json = json.load(inputfile)

    for examiner in input_json:
        for report in examiner["reports"]:
            semester, year = date2semester(report["date"])

            report["year"] = year
            report["semester"] = semester
            report["subject"] = norm_subject_name(report["subject"])
            report["examiners"] = examiner["examiner"]
            report["file"] = None

            filepath = f"./download/{examiner['id']}/{report['id']}.pdf"

            report.pop("date", None)
            report.pop("id", None)

            upload_report(JWT_TOKEN, report, filepath)
            # print(json.dumps(report))


if __name__ == "__main__":
    main()
