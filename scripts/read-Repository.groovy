import groovy.json.JsonOutput
import groovy.json.JsonSlurper

if (args != "") {

  def configurationMap = retrieveRepositoryConfiguration(args)
  return JsonOutput.toJson(configurationMap)
}

def retrieveRepositoryConfiguration(String repositoryId) {

  def result = new HashMap<String, String>()

  def repository = repository.getRepositoryManager().get(repositoryId)

  try{
    //add actual configuration of repository
    result.put ("attributes", repository.getConfiguration().getAttributes())
  }
  catch (Exception e){
    if (e instanceof NullPointerException){
      return "404: no repository found"
    }
    return e
  }
  //add what Nexus calls meta-data
  result.put ("id", repository.getName())
  result.put ("format", repository.getFormat().value)
  result.put ("type", repository.getType().getValue())

  return result
}

