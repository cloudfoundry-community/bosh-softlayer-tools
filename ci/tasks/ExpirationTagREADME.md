#Expiration Tags for Swift Object Storage

Swift Object Storage has expiring object support which can be used to schedule automatic deletion of objects within a day of a certain date. This feature can be set in a couple ways. 

##Setting the Tag through Swift API:

The steps to authenticate are outlined in the [Managing SoftLayer Object Storage Through REST APIs](http://sldn.softlayer.com/blog/waelriac/Managing-SoftLayer-Object-Storage-Through-REST-APIs) article under the Managaing the Object Storage section.

Once you have both X-Auth-Token and the X-Storage-URL, the following API calls can be made. In the meta data for the the objects, there is a value for `X-Delete-At`. This value can be set with a `POST` request and one of the following headers: `X-Delete-At` or `X-Delete-After`. The first header, `X-Delete-At` requires an integer value that is the certain date, in UNIX Epoch timestamp format, when the object will be removed. The `X-Delete-After` requires an integer value that is the number of seconds after which the object will be removed. The value given through `X-Delete-After` is internally stored in `X-Delete-At` meta deta value. Using the authentication values, the following post request can be made to set the value:

```
curl -i $STORAGE_URL/<container>/<object_name> -X POST -H "X-Auth-Token: $AUTH-TOKEN" -H "X-Delete-After: 1209600"
```
A successful response will be:
```
HTTP/1.1 202 Accepted
Content-Length: 76
Content-Type: text/html; charset=UTF-8
X-Trans-Id: txb5fb5c91ba1f4f37bb648-0052d84b3f
Date: Thu, 16 Jan 2014 21:12:31 GMT

<html><h1>Accepted</h1><p>The request is accepted for processing.</p></html>
```
If the call is being made to remove the tag, substitute `"X-Delete-After: "` in the above command.
For more information about the different requests possible and to change other meta data values look at the [OpenStack Object Storage Documentation](http://developer.openstack.org/api-ref-objectstorage-v1.html).

##Setting the Tag through SoftLayer Customer Portal:

Another way to see the value of the tag, or to set it is to go to the SoftLayer Customer Portal. After navigating to the corresponding container, click on the object to be edited. In the properties for the object, seen on the right hand side of the screen, look for an `Expires At:` label with a box next to it which may or may not have a date in it. In there pick the date in which the object should be deleted or remove the expiration date by clearing the box. 
