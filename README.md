# Supergoon Dynamic DNS updater
#### Updates DNS records hourly incase of home DNS changes
- This should be run in a cronjob in k8s or something
## Inputs
- Need a DNS_ZONES ENV variable as config map
-- This should be a space separated a records in the zone
- Need a ZONE_ID ENV variable as config map
-- This should be your hosted zone ID to modify
- Also needs AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY for an IAM user who can update records

### Example env variables that you would pass in via swarm/k8s
#### Update the zones base and hello
- DNS_ZONES="base hello"
#### Your public hosted zone id
- ENV ZONE_ID="1234567891234567890"
#### Your AWS IAM USER creds
- ENV AWS_ACCESS_KEY_ID=ACCESSKEYID
- ENV AWS_SECRET_ACCESS_KEY=ACCESSKEYSECRET