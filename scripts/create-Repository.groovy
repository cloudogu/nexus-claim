import groovy.json.JsonSlurper
import org.sonatype.nexus.repository.config.Configuration

class Repository {
  Map<String, Map<String, Object>> properties = new HashMap<String, Object>()
}

if (args != "") {

  def rep = convertJsonFileToRepo(args)
  def output = createRepository(rep)

  return output

}

def createRepository(Repository repo) {

  def conf = createConfiguration(repo)

  try {
    repository.createRepository(conf)
  }
  catch (Exception e) {
    return e
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


def createConfiguration(Repository repo) {


  def name = getName(repo)
  def recipeName = getRecipeName(repo)
  def online = getOnline(repo)
  def attributes = repo.properties.get("attributes")

  attributes = configureAttributes(attributes)

  att = setDefaultValuesOnNull(attributes)

  Configuration conf = new Configuration(
    repositoryName: name,
    recipeName: recipeName,
    online: online,
    attributes: att
  )


  return conf
}

def setDefaultValuesOnNull(Map<String, Object> attributes) {

  Map<String, Object> result = attributes

  for (att in attributes) {

    if (att.key.equals("maven")) {

      HashMap<String, Object> mavenEntry = att.value
      mavenEntry = replaceMavenEntryIfNull(mavenEntry)
      result.put(att.key, mavenEntry)
    }
  }
  return result
}


def configureAttributes(Map<String, Object> attributes) {


  def newAttributes = attributes

  for (att in newAttributes) {

    if (att.key.equals("httpclient")) {

      HashMap<String, Object> mapValue = att.value.get(0)

      mapValue = configureAttributes(mapValue)

      att.setValue(mapValue)
      newAttributes.put(att.key, att.value)

    } else {
      newAttributes.put(att.key, att.value.get(0))

    }

  }
  return newAttributes
}

def isMapEmpty(Map <String, Object> entry ){

  if (entry.size() == 0) {
    return true
  }

  entry.each {
    if (!it.value.equals("")){
      return false
    }
  }

  return true

}

def replaceMavenEntryIfNull(Map<String, Object> entry) {

  if (isMapEmpty(entry)) {
    entry = getDefaultMavenEntry()
  }
  return entry

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
