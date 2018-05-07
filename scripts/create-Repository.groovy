import groovy.json.JsonSlurper
import org.sonatype.nexus.blobstore.api.BlobStoreManager
import org.sonatype.nexus.repository.config.Configuration
import org.sonatype.nexus.repository.storage.WritePolicy

/*
  WORK IN PROGRESS
 */
class Repository {
  Map<String, Map<String, Object>> properties = new HashMap<String, Object>()
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

  def conf = createHostedConfiguration(repo)

  try {
    repository.createRepository(conf)
  }
  catch (Exception e){
    return e
  }
  return "successfully created " + getName(repo)
}


def createHostedConfiguration(Repository repo){

  def name = getName(repo)
  def recipeName = getRecipeName(repo)
  def online = getOnline(repo)
  def attributes = repo.properties.get("attributes")

  if(recipeName.contains("proxy")){

    HashMap<String,Object> httpClient = attributes.get("httpclient")
    def connection = httpClient.get("connection").get(0)
    httpClient.put("connection",connection)


    attributes.put("proxy",attributes.get("proxy").get(0))
    attributes.put("negativeCache",attributes.get("negativeCache").get(0))
    attributes.put("httpclient",httpClient)

    attributes.put("storage", attributes.get("storage").get(0))
  }

  else if (recipeName.contains("hosted")){
    attributes.put("storage", attributes.get("storage").get(0))
  }

  if (recipeName.contains("maven")){

    attributes.put("maven", attributes.get("maven").get(0))
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

if (args != "") {

  def rep = convertJsonFileToRepo(args)
  def newRepo = createRepository(rep)

  return newRepo

}
