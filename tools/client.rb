require 'json_api_client'
require 'fileutils'
require 'net/http'

module Client
  # this is an "abstract" base class that
  class Base < JsonApiClient::Resource
    # set the api base url in an abstract base class
    self.site = 'https://moozean.mathphys.stura.uni-heidelberg.de/api/'
  end
  Base.connection do |connection|
    # set the HTTPBasicAuth Password
    connection.use Faraday::Request::BasicAuthentication, ENV['MOOZEAN_USERNAME'], ENV['MOOZEAN_PASSWORD']
    # log responses
    # connection.use Faraday::Response::Logger
  end
  class Subject < Base
  end
  class Typ < Base
  end
  class Examinator < Base
    has_many :reports
  end
  class IsAbout < Base
    has_one :modul
    has_one :report
  end
  class Report < Base
   has_many :moduls
   has_many :is_abouts
   has_many :examinators
   has_one :typ
   has_one :subject
  end
  class Modul < Base
  end
  class Person < Base
  end
  class Modul < Base
  end
end
typ = Client::Typ.where(name:"Klausur").first
examinators = Client::Examinator.all
read = true
profs = []
while read do
  examinators.each do |examinator|
    profs.push(examinator)
  end
  examinators = examinators.pages.next
  read = !(examinators.nil?)
end
result = []
profs.each do |prof|
  read = true
  # We created a special folder "Digitale Altklausurenausleihe" with id 102 that shows that a folder is allowed to be in the public exam interface
  reports = Client::Report.where("typ"=>typ.id, 'examinators'=> [prof.id], "folderseries"=>102).all
  reps = []
  while read do
    reports.each do |report|
      unless report.moduls.first.nil?
        reps.push({"date"=>report["examination-at"],"moduleName"=>report.moduls.first.name,"subject"=>report.subject["name"],"id"=>report["id"]})
      end
    end
    reports = reports.pages.next
    read =!(reports.nil?)
  end
  print(reports)
  unless reps.empty?
    result.push({"examiner"=>"#{prof.givenname} #{prof.surname}","id"=>prof.id,"reports"=>reps})
    prof_folder = "download/#{prof.id}"
    FileUtils.mkdir_p prof_folder
    reps.each do |report|
      dowload_destination = "#{prof_folder}/#{report['id']}.pdf"
      if not File.file?(dowload_destination)
        uri = URI("https://moozean.mathphys.info/api/download/#{report['id']}/3168")
        Net::HTTP.start(uri.host, uri.port,
          :use_ssl => uri.scheme == 'https',
          :verify_mode => OpenSSL::SSL::VERIFY_NONE) do |http|

          request = Net::HTTP::Get.new uri.request_uri
          request.basic_auth ENV['MOOZEAN_USERNAME'], ENV['MOOZEAN_PASSWORD']

          response = http.request request # Net::HTTPResponse object
          open(dowload_destination, "wb") do |file|
            file.write(response.body)
          end
        end
        print("Downloaded '#{dowload_destination}'!\n")
      else
        print("'#{dowload_destination}' already exists, not downloading!\n")
      end
    end
  end
end
outfile_name = "reports_raw.json"
File.write(outfile_name,result.to_json)
print("Wrote output to \"#{outfile_name}\"!\n")
