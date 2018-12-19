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
  catch (Exception e){
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

def configureAttributes(Map<String, Object> properties){

  for (pro in properties){

    if (pro.key.equals("httpclient")){

      HashMap<String,Object> mapValue = pro.value.get(0)

      mapValue = configureAttributes(mapValue)

      pro.setValue(mapValue)
      properties.put(pro.key, pro.value)

    }

    else {
      properties.put(pro.key, pro.value.get(0))

    }

  }
  return properties
}

def createConfiguration(Repository repo){

  def name = getName(repo)
  def recipeName = getRecipeName(repo)
  def online = getOnline(repo)
  def attributes = repo.properties.get("attributes")


  attributes = configureAttributes(attributes)

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

def getRecipeName(Repository repo){
  String recipeName = repo.getProperties().get("recipeName")
  return recipeName
}

def getOnline(Repository repo){
  String online = repo.getProperties().get("online")
  return online
}
