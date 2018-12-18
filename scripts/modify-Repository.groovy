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
    repository.getRepositoryManager().get(repositoryID).start()
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


  if(recipeName.contains("proxy")){
    attributes = configureProxyAttributes(attributes,recipeName)

  }

  else if (recipeName.contains("group")){
    attributes = configureGroupAttributes(attributes,recipeName)

  }
  else if (recipeName.contains("hosted")){
    attributes = configureHostedAttributes(attributes,recipeName)

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

def getOnline(Repository repo){
  String online = repo.getProperties().get("online")
  return online
}

def getRecipeName(Repository repo){
  String recipeName = repo.getProperties().get("recipeName")
  return recipeName
}

def configureGroupAttributes(Object attribute,recipeName){

  def attributes = attribute
  attributes.put("storage", attributes.get("storage"))
  attributes.put("group",attributes.get("group"))
  if (recipeName.contains("maven")){
    attributes.put("maven", attributes.get("maven"))
  } else if (recipeName.contains("docker")){
    attributes.put("docker", attributes.get("docker"))
  }
  return attributes
}

def configureHostedAttributes(Object attribute, String recipeName){

  def attributes = attribute
  attributes.put("storage", attributes.get("storage"))
  if (recipeName.contains("maven")){
    attributes.put("maven", attributes.get("maven"))
  } else if (recipeName.contains("docker")){
    attributes.put("docker", attributes.get("docker"))
  }

  return attributes
}

def configureProxyAttributes(Object attribute, String recipeName){

  def attributes = attribute
  HashMap<String,Object> httpClient = attributes.get("httpclient")
  def connection = httpClient.get("connection")
  httpClient.put("connection",connection)


  attributes.put("proxy",attributes.get("proxy"))
  attributes.put("negativeCache",attributes.get("negativeCache"))
  attributes.put("httpclient",httpClient)
  attributes.put("storage", attributes.get("storage"))

  if (recipeName.contains("maven")){

    attributes.put("maven", attributes.get("maven"))


  } else if (recipeName.contains("docker")){
    attributes.put("docker", attributes.get("docker"))
    attributes.put("dockerProxy", attributes.get("dockerProxy"))
  }

  return attributes
}
