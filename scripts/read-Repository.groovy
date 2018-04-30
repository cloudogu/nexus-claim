import groovy.json.JsonOutput
import groovy.json.JsonSlurper


//../groovy/nexus-scripting execute --payload test22 read-Repository.groovy

if (args != "") {
  // String id = parseRepositoryId(args)

  // return id


  def configurationMap = retrieveRepositoryConfiguration(args)

  return JsonOutput.toJson(configurationMap)
}

private String parseRepositoryId(String[] input) {


  def jsonData = new JsonSlurper().parseText(input[0])

  return jsonData.id
}

private Map<String, String> retrieveRepositoryConfiguration(String repositoryId) {

  def result = new HashMap<String, String>()

  def repository = repository.getRepositoryManager().get(repositoryId)

  //add actual configuration of repository
  result.put ("attributes", repository.getConfiguration().getAttributes())

  //add what Nexus calls meta-data
  result.put ("id", repository.getName())
  result.put ("format", repository.getFormat().value)
  result.put ("url", repository.getUrl())
  result.put ("type", repository.getType().getValue())

  return result
}

