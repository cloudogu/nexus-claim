import groovy.json.JsonSlurper
import org.sonatype.nexus.repository.config.Configuration
import java.util.*

class Repository {
  Map<String, Map<String, Object>> properties = new HashMap<String, Map<String, Object>>()
}

if (args != "") {
  def rep = convertJsonFileToRepo(args)
  def output = createRepository(rep)

  return output
}

def convertJsonFileToRepo(String jsonData) {
  def inputJson = new JsonSlurper().parseText(jsonData)
  Repository repo = new Repository()

  inputJson.each {
    repo.properties.put(it.key, it.value)
  }

  return repo
}

def createRepository(Repository repo) {
  def conf = createConfiguration(repo)

  try {
    repository.createRepository(conf)
  }
  catch (Exception e) {
    if (e.getMessage().contains("Failed to initialize facets")) {
      return new RuntimeException("Could not create repository. Does it already exist? Please check the logs. " + e.toString(), e)
    }
    return e
  }

  return null
}

def createConfiguration(Repository repo) {
  def name = getName(repo)
  def recipeName = getRecipeName(repo)
  def online = getOnline(repo)
  Map<String, Object> attributes = repo.properties.get("attributes")

  Configuration conf = new Configuration(
    repositoryName: name,
    recipeName: recipeName,
    online: online,
    attributes: attributes
  )

  return conf
}


def getName(Repository repo) {
  String name = repo.getProperties().get("repositoryName")
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
