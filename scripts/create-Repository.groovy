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

  attributes = flattenSingleObjectLists(attributes)

  att = possiblySanitizeMavenDefaults(attributes)

  Configuration conf = new Configuration(
    repositoryName: name,
    recipeName: recipeName,
    online: online,
    attributes: att
  )

  return conf
}

/* recursively converts map entries which are lists with exactly one object */
def flattenSingleObjectLists(Map<String, Object> attributes) {
  def attributesCopy = new HashMap<Map, Object>()

  for (Map.Entry entry : attributes) {
    def key = entry.getKey()
    def val = entry.getValue()

    def mapValue
    if (val instanceof List) {
      if (val.size == 1 && val.get(0) instanceof Map) {
        Map<String, Object> embeddedList = val.get(0)
        mapValue = flattenSingleObjectLists(embeddedList)
      } else {
        mapValue = val
      }
    } else {
      mapValue = val
    }
    attributesCopy.put(key, mapValue)
  }

  return attributesCopy
}

/* adds maven specific policies if they are unconfigured */
def possiblySanitizeMavenDefaults(Map<String, Object> attributes) {
  Map<String, Object> result = new HashMap<String, Object>(attributes)

  for (att in attributes) {
    if (att.key.equals("maven")) {
      HashMap<String, Object> mavenEntry = att.value

      mavenEntry = replaceMavenEntryIfNull(mavenEntry)
      result.put(att.key, mavenEntry)
    }
  }

  return result
}

def replaceMavenEntryIfNull(Map<String, Object> entry) {
  if (isMapEmpty(entry)) {
    entry = getDefaultMavenEntry()
  }

  return entry
}

def isMapEmpty(Map<String, Object> entry) {
  if (entry.size() == 0) {
    return true
  }

  entry.each {
    if (!it.value.equals("")) {
      return false
    }
  }

  return true
}

def getDefaultMavenEntry() {
  Map<String, String> mavenEntry = new HashMap<>()

  mavenEntry.put("layoutPolicy", "STRICT")
  mavenEntry.put("versionPolicy", "RELEASE")

  return mavenEntry
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
