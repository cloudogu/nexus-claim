import groovy.json.JsonOutput
import groovy.json.JsonSlurper

if (args != "") {

  def configurationMap = retrieveRepositoryConfiguration(args)

  return JsonOutput.toJson(configurationMap)
}

private Map<String, String> retrieveRepositoryConfiguration(String repositoryId) {

  def result = new HashMap<String, String>()

  def repository = repository.getRepositoryManager().get(repositoryId)

  //add actual configuration of repository
  result.put ("attributes", repository.getConfiguration().getAttributes())

  //add what Nexus calls meta-data
  result.put ("id", repository.getName())
  result.put ("format", repository.getFormat().value)
  result.put ("type", repository.getType().getValue())

  return result
}

