// Code generated by go generate; DO NOT EDIT.
package infrastructure

const CREATE_REPOSITORY = `import groovy.json.JsonSlurper
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
`
const DELETE_REPOSITORY = `def deleteRepository(String name) {

  try {
    repository.getRepositoryManager().delete(name)
  }
  catch (Exception e){
    return e
  }
  return null
}

if (args != "") {

  def output = deleteRepository(args)

  return output
}

`
const MODIFY_REPOSITORY = `import org.sonatype.nexus.common.stateguard.InvalidStateException
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
`
const READ_REPOSITORY = `import groovy.json.JsonOutput
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

`
