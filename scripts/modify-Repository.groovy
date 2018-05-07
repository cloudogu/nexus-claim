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

  return "successfully modified " + repositoryID
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
  def online = getOnline(repo)
  def attributes = repo.properties.get("attributes")

  if(recipeName.contains("proxy")){
    attributes = putProxyAttribute(attributes,recipeName)
  }

  else if (recipeName.contains("group")){
    attributes = putGroupAttribute(attributes)

  }
  else if (recipeName.contains("hosted")){
    attributes = putHostedAttribute(attributes,recipeName)
  }

  Configuration conf = new Configuration(
    repositoryName: name,
    recipeName: recipeName,
    online: online,
    attributes: attributes
  )


  return conf
}

def getName(Repository repo){
  String name = repo.getProperties().get("name")
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

def putGroupAttribute(Object attribute){

  def attributes = attribute
  attributes.put("storage", attributes.get("storage").get(0))
  attributes.put("group",attributes.get("group").get(0))
  return attributes
}

def putHostedAttribute(Object attribute, String recipeName){

  def attributes = attribute
  attributes.put("storage", attributes.get("storage").get(0))
  if (recipeName.contains("maven")){
    attributes.put("maven", attributes.get("maven").get(0))
  }

  return attributes
}

def putProxyAttribute(Object attribute,String recipeName){

  def attributes = attribute
  HashMap<String,Object> httpClient = attributes.get("httpclient")
  def connection = httpClient.get("connection").get(0)
  httpClient.put("connection",connection)

  attributes.put("proxy",attributes.get("proxy").get(0))
  attributes.put("negativeCache",attributes.get("negativeCache").get(0))
  attributes.put("httpclient",httpClient)
  attributes.put("storage", attributes.get("storage").get(0))

  if (recipeName.contains("maven")){
    attributes.put("maven", attributes.get("maven").get(0))
  }

  return attributes
}
