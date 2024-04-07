# GoProfileFromSeed 🕵️🏾‍♂️
## A Go library to generate user profiles based on a seed.


A `Profile` is a structure with the following fields

-   **First Name**
-   **Last Name**
-   **Username**
-  **Email**
-  **Region**
-  **Address**
-  **Seed**

**You can also generate a unique profile image for each profile.**

A `seed` is a 5 character string. Each character can a number or a lower case or upper case English letter.

## How it works 🛠️
The following image shows how each field is determined :

![enter image description here](https://i.ibb.co/d2JZYnv/seed-1.png)

**Region**:  numbers, upper case letters and lower case letters are divided into groups and if the integer representation of the byte (n) of a character is within that group that region is assigned. Same for gender.

**First names starting letter** will be either,
the first non-numerical character after the second character of the seed or if no non-numerical character is present after the second character of the seed,  then the letter corresponding to the integer representation of the byte of the 3rd number + 17.

**Last names starting letter** will be either,
the second non-numerical character after the second character of the seed or if only one non-numerical character is present after the second character of the seed,  then the letter corresponding to the integer representation of the byte of the 4th number + 17.

**Offset** : The program will take the integer representation of the byte (n) of a character and will assign the item at the n th line in the data file belonging to the respective field.  In the case of two letters, it will multiply those two numbers and get the n.

**Username and email** is generated using templates determined in the common_templates directory in the data.


## Profile Image 
Each profile image is a 250 * 250 pixel image. This square is divided into 25, 50*50 squares and according to the image below, five or fewer squares are painted based on each character of the seed starting from the middle square.

![enter image description here](https://i.ibb.co/Y8F2DXF/profile-img.png)
Examples:
![enter image description here](https://i.ibb.co/D83r0rj/examples.png)
## Data files

Currently, the data files include common names and addresses of these regions:

- Sri Lanka 🇱🇰
- USA 🇺🇸
- UK 🇬🇧
- Australia 🇦🇺

Information on how those data is gathered can be found in info.txt files inside each folder.

## Documentation
https://pkg.go.dev/github.com/Tharusha-dev/GoProfileFromSeed
