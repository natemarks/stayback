# stayback

stayback is a command line tool used for storing/retrieving sensitive linux desktop data in various AWS services.

sbs3  is for storing big bundles of directoris in s3, optionally encrypted
sbsecrets is for packing a small local file into secrets manager
sbkp is for packing a locak ssh key into EC2 keypair. this is a 1-way operation