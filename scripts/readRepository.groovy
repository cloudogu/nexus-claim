import groovy.json.JsonOutput
import groovy.json.JsonSlurper

if (args != "") {
  String id = parseRepositoryId(args)
  def configurationMap = retrieveRepositoryConfiguration(id)

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
  result.put ("name", repository.getName())
  result.put ("format", repository.getFormat())
  result.put ("url", repository.getUrl())
  result.put ("type", repository.getType().getValue())

  return result
}

