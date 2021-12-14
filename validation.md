the External Amazon DB HOSTNAME: ate-restaurant-services.ciy4rgtokd8m.us-east-1.rds.amazonaws.com

# Schema

Changed the File names from:    
tbl_user.sql -> 0_tbl_user.sql  
tbl_query.sql -> 1_tbl_query.sql

# DataBase

Updated The Database Tables:

1.) Added the user_type Table. That will contain the 4 types:

Customer, Restaurant, Driver, Admin

2.) users: Added All the multiple attributes to define the needs of the Kitchen.

cuisine          
status           
everyday        
profileImageId   
shopLogoId      
shopBannerId     
isVeg            
mealService      
partyCatering    
deliveryTakeAway    
delivery         
freeDelivery     
offerType        
offer            
offerAmount      
maxDeliveryTime  
description      
location         
locLongitude     
locLatitude

The "type" attribute will determine which type the user is.

The Attributes profileImageId, shopLogoId, shopBannerId will contain the Images related to the User.

3.) Added The open_time Table related to Opening Days / Time of the Restaurant

This tables contains information regarding each restaurant opening days:
DayName, From, to times.

4.) Added to the product table the Attribute:

productImage

this attribute will contain the image url.

5. Added The Tables: tables, reservations for the Table Reservation Service using QR Technologies.

# Menu Service

Updated The menu-service.proto:

1-) adding File Upload for Product Image.

2-) Regenerated the Swagger / GRPC Files.

3-) Updated the Redefinition of the GRPC Methods and The SQL Requests.

4-) File Upload is done by Uploading the FileName + FileBytes.

5-) Removed some code unrelated to menu service due to fast copy/past.

6-) Added TLS to gRPC.

7-) TLS was not added to http.Client due to some errors. (need to be much further looked into.)

# Setting Service

1-) Definition of the Proto files with full attributes and methods.

2-) Generation of the Swagger / GRPC Files.

3-) Redefinition of the GRPC Methods + SQL Requests.

4-) File Upload is done by Uploading the FileName + FileBytes for the 3 Images.

5-) Added TLS to gRPC.

6-) TLS was not added to http.Client due to some errors. (need to be much further looked into.)

# Barcode Service

1-) Fixed some code. Removed some unrelated Code and still in progress.

# Reservation Service

1-) Some Code fix. but no big progress so far.

# Nginx

1-) Added some code for the Setting, Addon Services.
