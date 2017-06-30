# AWS-RDS-LOG-TAILER

A tool to fetch aws rds log continuously, base on `AWS SKD Go` library.

# Prerequisite

You need run this binary in AWS EC2 instance who has IAM Role to access AWS RDS 

# Usage

Just run `bin/amd64/aws-rds-log-tailer` with proper options.

```bash
./aws-rds-log-tailer -dbID rds-instance-identifier -out postgres.log
```

# Development

This project is also docker base project, you can use:

* `make`, to build new binary
* `make test`, to run unit test
* `make enter`, to enter the container to do debug

And it needs aws credential to run the aws sdk api, and did not implement this part yet, instead
using aws role to do this job, you can copy the binary(`bin/amd64/aws-rds-log-tailer`) to ec2 instance(like by s3), and do the test

# Reference

http://docs.aws.amazon.com/cli/latest/reference/rds/download-db-log-file-portion.html

# TODO

- [ ] refactor with goroutine 
- [ ] add yml conf file to config different databases and output files
