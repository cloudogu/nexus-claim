import org.sonatype.nexus.common.stateguard.InvalidStateException
import org.sonatype.nexus.repository.config.Configuration
import groovy.json.JsonSlurper

class Repository {
  Map<String, Map<String, Object>> properties = new HashMap<String, Object>()
}

if (args != "") {

  def repo = convertJsonFileToRepo(args)
  def name = getName(repo)
  def conf = createConfiguration(repo)

  def output = modifyRepository(name, conf)
  return output
}

def modifyRepository(String repositoryID, Configuration configuration) {

  repository.getRepositoryManager().get(repositoryID).stop()

  try {
    repository.getRepositoryManager().get(repositoryID).update(configuration)
  }
  catch (Exception e) {
    return e
  }
  finally {
    repository.getRepositoryManager().get(repositoryID).start()
  }

  return null
}

def convertJsonFileToRepo(String jsonData) {

  def inputJson = new JsonSlurper().parseText(jsonData)
  Repository repo = new Repository()
  inputJson.each {
    repo.properties.put(it.key, it.value)
  }

  return repo
}

def createConfiguration(Repository repo){

  def name = getName(repo)
  def recipeName = getRecipeName(repo)
  def attributes = repo.properties.get("attributes")
  def online = getOnline(repo)

  Configuration conf = new Configuration(
    repositoryName: name,
    recipeName: recipeName,
    online: online,
    attributes: attributes
  )

  return conf
}

def getName(Repository repo){
  String name = repo.getProperties().get("repositoryName")
  return name
}

def getOnline(Repository repo){
  String online = repo.getProperties().get("online")
  return online
}

def getRecipeName(Repository repo){
  String recipeName = repo.getProperties().get("recipeName")
  return recipeName
}
