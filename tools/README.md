# How to:

1. `$ bundle install`
2. Set the enviroment Variables `MOOZEAN_USERNAME` and `MOOZEAN_PASSWORD`
3. `$ bundle exec ruby client.rb`
4. Reset the Moozean
    - Reset the exam bucket:
      ```
      $ mc move -r "mathphysinfo/altklausuren" "mathphysinfo/altklausuren-archive/altklausuren-$(date +%Y-%m-%d-%H_%M)"
      ```
    - Reset the db:
      ```
      cherry:/opt/docker/altklausur-ausleihe# docker-compose -f docker-compose-production.yml down --volumes
      cherry:/opt/docker/altklausur-ausleihe# docker-compose -f docker-compose-production.yml up -d
      cherry:/opt/docker/altklausur-ausleihe# docker-compose -f docker-compose-production.yml logs -f
      ```
5. Redeem an Admintoken at https://altklausuren.mathphys.info/adminlogin
6. `$ python uploadscript.py`
