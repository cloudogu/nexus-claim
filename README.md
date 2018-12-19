# Nexus Claim

Define your [Sonatype Nexus](http://www.sonatype.org/nexus/) repository structure as code.


## Example uses of nexus-claim on nexus 3
```bash
$ nexus-claim plan -i resources/nexus3/nexus3-initial-example.hcl -o resources/nexus3/nexus3-initial-example.json
$ nexus-claim apply -i resources/nexus3/nexus3-initial-example.json
```

## Things to consider for .hcl of nexus 3
- `repository "xxx"` should be the same as `repositoryName = "xxx"`
- on maven2-hosted- and maven2-proxy-repository there must be a maven sector in addition(with versionPolicy and writePolicy) 
 
## Example .hcl for nexus 3

can be found [here](https://github.com/cloudogu/nexus-claim/blob/feature/3_fix_create_docker_repo/resources/nexus3/nexus3-initial-example.hcl)
 

## Example .hcl for nexus 2

can be found [here](https://github.com/cloudogu/nexus-claim/blob/feature/3_fix_create_docker_repo/resources/nexus-initial-example.hcl)




