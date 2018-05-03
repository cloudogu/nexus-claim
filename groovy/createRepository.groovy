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

// Not needed
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

    return "test"

    def inputJson = new JsonSlurper().parseText(jsonData)
    Repository repo = new Repository()
    inputJson.each {
        repo.properties.put(it.key, it.value)
    }

    return repo
}

def createRepository(Repository repo) {

    Configuration conf = new Configuration()


    if (getRecipeName(repo).contains("hosted")){

        conf = createHostedConfiguration(repo)
    }
    if (getRecipeName(repo).contains("proxy")){

    }

    repository.createRepository(conf)

    return "successful"
}

def createProxyConfiguration(Repository repo){

}

def createHostedConfiguration(Repository repo){

    def name = getName(repo)
    def recipeName = getRecipeName(repo)
    def attributes = repo.properties.get("attributes")


    Configuration conf = new Configuration(
            repositoryName: name,
            recipeName: recipeName,
            online: true,
            attributes: attributes
    )

    if (recipeName.contains("maven")){
        conf.attributes.maven = repo.getProperties().get("attributes").get("maven")
    }

    return conf

}

def getName(Repository repo){
    String name = repo.getProperties().get("id")
    return name
}

def getRecipeName(Repository repo){
    String recipeName = repo.getProperties().get("format") + "-" + repo.getProperties().get("type")
    return recipeName
}


if (args != "") {

    def rep = convertJsonFileToRepo(args)

    def newRepo = createRepository(rep)

    return newRepo

}
