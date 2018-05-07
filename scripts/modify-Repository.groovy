import org.sonatype.nexus.common.stateguard.InvalidStateException
import org.sonatype.nexus.repository.config.Configuration
import groovy.json.JsonSlurper

if (args != "") {


  def repo = convertJsonFileToRepo(args)
  def name = getName(repo)

  def conf = createHostedConfiguration(repo)

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

  return "successfully modified " + repositoryID
}


def createHostedConfiguration(Repository repo) {

  def name = getName(repo)
  def recipeName = getRecipeName(repo)
  def online = getOnline(repo)

  def attributes = repo.properties.get("attributes")
  attributes.put("storage", attributes.get("storage").get(0))

  Configuration conf = new Configuration(
    repositoryName: name,
    recipeName: recipeName,
    online: online,
    attributes: attributes
  )

  if (recipeName.contains("maven")) {
    conf.attributes.maven = repo.getProperties().get("attributes").get("maven")
  }

  return conf

}


class Repository {
  Map<String, Map<String, Object>> properties = new HashMap<String, Object>()
}

def getName(Repository repo) {
  String name = repo.getProperties().get("name")
  return name
}

def getRecipeName(Repository repo) {
  String recipeName = repo.getProperties().get("recipeName")
  return recipeName
}

def getOnline(Repository repo) {
  String online = repo.getProperties().get("online")
  return online
}

def convertJsonFileToRepo(String jsonData) {

  def inputJson = new JsonSlurper().parseText(jsonData)
  Repository repo = new Repository()
  inputJson.each {
    repo.properties.put(it.key, it.value)
  }

  return repo
}

def isRepositoryManagerStarted() {

  try {
    repository.getRepositoryManager().start()
    return true
  }
  catch (Exception e) {
    return false
  }
}