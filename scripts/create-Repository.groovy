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

def getDefaultAttributeValues(String type, String remoteURL){

  Map<String, Map<String, Object>> defaultAttributeValues = new HashMap<String,Object>()

  if (type.equals("hosted")) {

    Map<String,Object> storage = new HashMap<String,Object>()
    storage.put("blobStoreName",BlobStoreManager.DEFAULT_BLOBSTORE_NAME)
    storage.put("writePolicy",WritePolicy.ALLOW)
    storage.put("strictContentTypeValidation",false)
    defaultAttributeValues.put("attributes",storage)
  }
  else if (type.equals("proxy")){
    Map<String,Object> httpClient = new HashMap<String,Object>()
    Map<String,Object> connection = new HashMap<String,Object>()
    connection.put("blocked",false)
    connection.put("autoBlock",true)
    httpClient.put("connection",connection)
    defaultAttributeValues.put("httpclient",httpClient)

    Map<String,Object> proxy = new HashMap<String,Object>()
    proxy.put("remoteUrl",remoteURL)
    proxy.put("contentMaxAge",1440)
    proxy.put("metadataMaxAge",1440)
    defaultAttributeValues.put("proxy",proxy)

    Map<String,Object> negativeCache = new HashMap<String,Object>()
    negativeCache.put("enabled",true)
    negativeCache.put("timeToLive",1440)
    defaultAttributeValues.put("negativeCache",negativeCache)

    Map<String,Object> storage = new HashMap<String,Object>()
    storage.put("blobStoreName",BlobStoreManager.DEFAULT_BLOBSTORE_NAME)
    storage.put("strictContentTypeValidation",false)
    defaultAttributeValues.put("attributes",storage)

  }
  else if(type.equals("group")){
    Map<String,Object> storage = new HashMap<String,Object>()
    storage.put("blobStoreName",BlobStoreManager.DEFAULT_BLOBSTORE_NAME)
    defaultAttributeValues.put("attributes",storage)

  }

  return defaultAttributeValues
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


  String name = repo.properties.get("name")
  repo.properties.remove("name")

  String recipeName = repo.properties.get("recipeName")
  repo.properties.remove("recipeName")

  String state = repo.properties.get("_state")
  repo.properties.remove("_state")



  Configuration conf = new Configuration(
    repositoryName: name,
    recipeName: recipeName,
    online: repo.properties.get("online"),
    attributes: [
      storage: [
        blobStoreName              : repo.properties.get("blobStoreName"),
        writePolicy                : writePolicy,
        strictContentTypeValidation: strictContentTypeValidation
      ] as Map
    ] as Map
  )

  conf.setAttributes()

  repository.createRepository(conf)

}


if (args != "") {

  def rep = convertJsonFileToRepo(args)
  createRepository(rep)

}
