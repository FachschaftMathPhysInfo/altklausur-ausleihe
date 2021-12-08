import json
from datetime import datetime
import os.path


def upload_report(JWT_TOKEN, report_metainformation, filepath):

    # Abort if any argument is not given
    if not os.path.isfile(filepath):
        print(f"Error, file \"{filepath}\" not found!")
        return

    os.system(
        f"./upload_one.sh {JWT_TOKEN} \'{json.dumps(report_metainformation)}\' {filepath}"
    )


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
    JWT_TOKEN = "Testtoken"

    # load the inputs
    with open("./reports_raw.json") as inputfile:
        input_json = json.load(inputfile)

    for examiner in input_json:
        for report in examiner["reports"]:
            report_metainformation = report
            semester, year = date2semester(report["date"])

            report_metainformation["year"] = year
            report_metainformation["semester"] = semester

            filepath = f"./download/{examiner['id']}/{report['id']}.pdf"

            report_metainformation.pop("date", None)
            report_metainformation.pop("id", None)

            upload_report(JWT_TOKEN, report_metainformation, filepath)
            # print(json.dumps(report_metainformation))


if __name__ == "__main__":
    main()
