---
title: "Review of my first Project"
url: ["blog/review-of-my-first-project"]
description: "A review of my first project."
date: "2020-04"
---

## Intro

In 2012 I (13 years old) started learning coding. I didn't read books or anything,
I mostly copied stuff from YouTube tutorials. This lead to my first bigger project.
A social network clone, which should be similar to Facebook.

In this blog post, I'll review my now 8 year old code.
Partly to amuse myself what kind of errors I made and
to maybe learn something from it.

The project has the following features:
* post timeline with like and comment functionality
* "real-time" updates of comments and like counts for posts
* chat between two users
* friend system
* user search
* profile page
* profile picture upload
* login and register with email confirmation

It looks like this:

![](/img/3-header.png)

Yes, most of the stuff is in German, but I'll rename variables
and strings in code examples to English.

## Code Style/Quality

Style guides define rules on how code should look like.
This is mostly done to enforce
consistency, maintainability and reduce common programming errors.
For PHP there are code styling guidelines defined in
[PSR-1](https://www.php-fig.org/psr/psr-1/) and [PSR-12](https://www.php-fig.org/psr/psr-12/).

My project didn't follow any of such guidelines
and it sure looks like garbage.
Bad code is the only consistency in that project :).
In newer projects I use linters like eslint and formatters
like prettier to ensure code quality and style.

Example of ugly code.
```php
$number = 2;
$query = mysql_query("SELECT * FROM post_comment WHERE post_id='".$post_id."' ORDER BY comment DESC");
	WHILE ($row = mysql_fetch_assoc($query)){
	if($zahl != 0) {
	$comment 	= $row["comment"];
	$name		= $row["name"];
	$date	 	= $row["date"];
		 if(strlen($comment)>=47)
		{
			$comment = substr($comment, 0, 47)."...";
		}
		$end .= "<hr><a href='/profile/".$name."'><b>"    .user_return($name)."</b></a><b>|</b> ".$comment;
		$number = $number - 1;
		}
	}
	echo $end;
```

### Unused Code

Having less code greatly increases code maintainability,
code complexity, reduces potential bugs, and overall simplicity of a project.
New team members not only have to understand used code but also the unused code
which doesn't add any benefit and is only a waste of time.

Without code, it is possible to write a secure and reliable application -> 
[github.com/kelseyhightower/nocode](https://github.com/kelseyhightower/nocode)

When parts of an application will be refactored. It is likely, 
that code will lose its callers and be unused.
That's okay as long as the code doesn't get into production.
Refactorings should be done inside a branch inside version control
and be reviewed for code quality before going into production.

### Reuse Code

Similar code should not be copied without reasoning.
Depending on the context, it maybe perfectly fine to copy things.
F.ex. in tests I prefer having all the relevant
test setup inside the actual test method and not extracting
it to another method.

In my project, there is a lot copied code. I think the main reason was,
that I really rarely used `return` in functions and
rather just `echo`'ed/printed them out.

A prime example for code that should be reused is the rendering of
posts inside the timeline. In there I created two functions one for
the timeline for the current user, and one for the profile page of a user.
The two functions are like 80 chars each and are pretty much the same,
only the SQL query for getting posts is a little different.

### File Structure

Structuring a project isn't easy.
While I extracted code to different files, I didn't really grouped it well.
There is a file in the project called `action/functions.php`, and like the name tells,
it contains nearly all business logic. Every function in that file, has a header.
It looks like this:
```php
///////////////////////////////////////////
///////////////////////////////////////////
///////////////See All Posts///////////////
///////////////////////////////////////////
///////////////////////////////////////////
function See_all_post($id) {
    // content
}
```
I think I knew what the problem was. But I tried to fix it the wrong way, with literally
only garbage comments. It also seems like, I didn't knew `ctrl+f` existed :D.
Anyway, the clear solution for that problem is properly structuring the project.
Each feature should get its own folder where views, business logic and routing should stay.
With this, finding function should be easy.

### Automated Tests

Obviously I didn't wrote tests. Why should I?
Most of the YouTube tutorials only specified how to write something,
not how to write good maintainable code. Tests verify that your code *really* works.
Here is a good summary why writing tests rocks!
[StackOverflow: Is Unit Testing worth the effort?](https://stackoverflow.com/a/67500/4244993)

In [gotify/server](https://github.com/gotify/server) 
I've written unit, integration and end to end tests from the start.
It is such a nice feeling to add a feature and still *know* that the rest still works as intended,
because the code is properly tested.

## Security

### SQL Injection

Basically every SQL statement in that project is vulnerable. Example, email verification:
```php
mysql_query("UPDATE user SET activated = '1' WHERE id='".$_GET['code']."'");
```

`$_GET['code']` comes from an query parameter.
The user can type anything in it. Like f.ex. malicious SQL like `irrelevant' OR 1 = 1 -- `.
This would lead to the following query:
```sql
UPDATE user SET activated = '1' WHERE id='irrelevant' OR 1 = 1 -- '
```
After executing this query, every email would be verified, without actually receiving the email.

I've found one query where I used `mysql_real_escape_string`.
```php
mysql_query("INSERT INTO post_comment ( post_id, name, comment, date) VALUES ('"
   .$post_id."','".$id."','".mysql_real_escape_string($message)."','"
   .$date."')")or die(mysql_error());
```
Escaping the message is good, but `$id` is also user supplied. I think I can safely say,
that nearly all inputs on that website are vulnerable.

SQL Injection is on place 1 on the [OWASP top ten security risks](https://owasp.org/www-project-top-ten/).
Without knowledge it is pretty easy to allow SQL injection and
with little knowledge it is also easy to prevent it.

The cure are prepared statements. These statements prevent the injection of sql.

In PHP prepared statements can be used via [PDO](https://www.php.net/manual/en/book.pdo.php).
The above `mysql_query` can be fixed like this.

```php
$pdo = new PDO('mysql:host=myhost;dbname=social', 'root', 'password');

$statement = $pdo->prepare("UPDATE user SET activated='1' WHERE id=?");
$statement->execute([$_GET['code']]);
```

### Cross-Site Request Forgery (CSRF)

In a CSRF attack, the attacker tries to let the victim submit a malicious web request without knowing it.
This could be to gain access to administration stuff or in my case add a new post to the timeline.

Here a vulnerable form:
```php
<form action="index.php" method="post">
    <textarea name="posttext"></textarea>
    <input type="submit" name="posten" value="Posten">
</form>
<?php
if ($_POST["posten"]) {
    $text = $_POST["posttext"];
    if ($text != "") {
        post_schreiben($text, $userid);
    } else echo "<p id='falsered'>Content my not be empty</p>";
}
?>
```

This following website submits the malious form instantly after visiting the page,
the user only have to navigate to this page,
and then the message `INJECT MESSAGE` will be posted on the timeline.
```html
<!DOCTYPE html>
<html>
<body>
    <form action="http://example.org/index.php" method="post">
        <input type="hidden" name="posttext" value="INJECT MESSAGE"/>
        <input type="hidden" name="posten" value="nah"/>
    </form>
    <script>
        document.forms[0].submit();
    </script>
</body>
</html>
```

![](/img/3-attack.gif)

Guarding against CSRF is somewhat more complicated.
One way is to add an anti forgery token to the form.
This token is generated on the server thus,
can't be send from the malious website because it has no knowledge of it.
After receiving a from post request from a client,
the server then validates the anti forgery token
and only then executes the request.

Pseudo code for a fix.

```php
<?php
$token = createNewCSRFToken()
?>
<form action="index.php" method="post">
    <textarea name="posttext"></textarea>
    <input type="hidden" name="csrf_token" value="<?=$token?>"
    <input type="submit" name="posten" value="Posten">
</form>
<?php
if ($_POST["posten"]) {
    if (checkIfCSRFTokenIsValid($_POST['csrf_token'])) {
        $text = $_POST["posttext"];
        if ($text != "") {
            write_post($text, $userid);
        } else echo "<p id='falsered'>Content my not be empty</p>";
    } else echo "<p id='falsered'>Something went wrong try again.</p>";
}
?>
```

### Cross-Site-Scripting (XSS)

Cross-Site-Scripting means the attacker injects HTML/JavaScript into the website.
The vulnerable part of my project is the comment section of posts.
The content of the actual post is secured, guess I partly knew what I was doing.

![](/img/3-cross-site-scripting.gif)

```php
if ($_GET["comment"] != "") {
    if (comment_post($userid, $_GET["comment"], $_GET["post_id"])) {
        echo "Success";
    } else echo "Something went wrong";
} else echo "Comment may not be empty";
```
To secure the script above, we need to sanitize `$_GET["comment"]`.
In this case `htmlentities` can be used.
This function converts all applicable characters to HTML entities.

F.ex. `<script>` will be converted to `&lt;script&gt;`.
This prevents injecting html into the content of the comment.

```php
if ($_GET["comment"] != "") {
    if (comment_post($userid, htmlentities($_GET["comment"]), $_GET["post_id"])) {
        echo "Success";
    } else echo "Something went wrong";
} else echo "Comment may not be empty";
```

### User Passwords

Passwords aren't easy. [This stackoverflow answer](https://security.stackexchange.com/a/31846)
summarizes why passwords need to be secured how to secure them properly.

TL;DR: Use a cryptographic hash functions, with a salt and make it slow (:.

My project uses md5 as hash function without any salt.
This is bad because it allows the attack with
[Rainbow Tables](https://en.wikipedia.org/wiki/Rainbow_table),
this is a table with precomputed hashes, thus allowing a quick lookup of a hash.

Adding a salt would still not be enough, because md5 is to fast and can be easily brute forced.

PHP provides an simple api for creating and validating secure passwords hashes with good defaults.

```php
$pw = "mypw"
$pwHash = password_hash($mypw, PASSWORD_DEFAULT);
// $pwHash = $2y$10$.vGA1O9wmRjrwAVXD98HNOgsNpDczlqm3Jq7KnEd1rVAGv3Fykk1a
// this hash contains the algorithm and other parameters
// the php default currently is BCrypt

if (password_verify($mypw, $pwHash)) {
    echo "SUCCESS";
}
```

### Hardcoded Passwords

An application internally uses passwords,
be it the database credentials or credentials for an online storage like aws.

Passwords like this should never be inserted plainly into the code
because it is pretty easy to commit them into version control like git
or make them browsable via a web server.
See [1% of CMS-Powered Sites Expose Their Database Passwords](https://feross.org/cmsploit/).

Passwords should either be stored in environment variables or 
inside config files outside of the web root.

In this project, the credentials of the mysql database was
stored plainly multiple times in different source files.
There had also different passwords, tho there is a lot of unused code inside the project :D.

```php
<?php
mysql_connect('mysql', 'root', 'password');
mysql_select_db('dbname');
?>
```

### Serverside Authorization

Users should only be able to access/edit resources which they have permission for.
F.ex. adding comments to a post from a user that isn't a friend should be prohibited.

Code without authorization check:
```php
if ($_GET["comment"] != "") {
    if (comment_post($userid, $_GET["comment"], $_GET["post_id"])) {
        echo "Success";
    } else echo "Something went wrong";
} else echo "Comment may not be empty";
```

Code with authorization:
```php
if ($_GET["comment"] != "") {
    if (postIsVisibleToUser($userid, $_GET["post_id"])) {
        if (comment_post($userid, $_GET["comment"], $_GET["post_id"])) {
            echo "Success";
        } else echo "Something went wrong";
    } else echo "Unauthorized";
} else echo "Comment may not be empty";
```

### User Input

User input should never be trusted, and always be validated.
This includes emails, dates, numbers and so on.

Not validating input can lead to bugs because some parts of the application
may depend on valid inputs. A good example for this, is the registration.
In there the user often adds an email address.
If the users provides an invalid email address, intentionally or not
the application should present the user an clear error message.
Otherwise, password reset or similar functionality may not work.

## WTF is this?

### Register

```php
function register($name, $npasswort, $email)
{
    $id1 = rand(1000000000, 9999999999);
    $id2 = rand(1000000000, 9999999999);
    $id3 = rand(1000000000, 9999999999);
    $id4 = rand(1000000000, 9999999999);
    db();
    $sql = mysql_query("SELECT * FROM user");
    $row = mysql_fetch_assoc($sql);
    if ($row['email'] != $email) {
        if ($row['id'] == $id1) {
            if ($row['id'] == $id2) {
                if ($row['id'] == $id3) {
                    if ($row['id'] == $id4) {

                    } else $id = $id4;
                } else $id = $id3;
            } else $id = $id2;
        } else $id = $id1;

        if ($id != "") {
            // insert user
        }
    }
}
```

Each user gets an unique id, this part tries to ensure that the id and email is unique.
And well, it only checks against the first column inside the table. I say, at least
one user will have a unique id. Having a duplicate email, is kinda bad, because it is used
for logging in. If the ID is not unique, 4 random generated ids will be checked. If every id exists
which is impossible, because only the first row will be checked, the script just does nothing.

### Classes where they shouldn't be

```php
class Login {
    protected $_email, $_password, $_result;

    public function __construct($email, $password) {
        $this->_email = $email;
        $this->_password = $password;
    }

    public function Login() {
        $db = new Database();
        if ($this->_email != "" or $this->_password != "") {
            if ($this->_email != "E-Mail" and $this->_password != "password") {
                if (login($this->_email, $this->_password)) {
                    // set user onto session
                    $this->_result = "Success";
                }
            } else $this->_result = "Fields may not be empty";
        } else $this->_result = "Fields may not be empty";
        $db->disconnect();
    }

    public function result() {
        return $this->_result;
    }
}
```
Usage:
```php
$Login = new Login($form_email, $form_password);
$Login->Login();
echo $Login->result();
```
This is just a pretty bad example of using classes.
Because it give no benefit but just bloats the code.
The whole class should just be an function.
```php
echo loginUser($form_email, $formPassword);
```

## Conclusion

As expected my social network delivered,
many security vulnerabilities, bad code quality and some WTF moments.
I enjoyed looking over it. It was a fun ride through the past :D.
Thank you 13 year old me, for saving this project on cloud storage
and even creating a sql dump from the database. Much appreciated.
