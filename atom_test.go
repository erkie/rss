package rss

import "testing"

func TestParseContentWithoutCDATA(t *testing.T) {
	doc := `
  <?xml version="1.0" encoding="UTF-8"?>
<feed xml:lang="en-US" xmlns="http://www.w3.org/2005/Atom" xmlns:activity="http://activitystrea.ms/spec/1.0/" xmlns:georss="http://www.georss.org/georss" xmlns:ostatus="http://ostatus.org/schema/1.0" xmlns:thr="http://purl.org/syndication/thread/1.0" xml:base="https://aaronparecki.com/">
    <generator uri="https://granary.io/">granary</generator>
    <id>https://aaronparecki.com/articles</id>
    <title>User feed for Aaron Parecki</title>
    <logo>https://aaronparecki.com/images/profile.jpg</logo>
    <updated>2019-02-25T09:58:24-06:00</updated>
    <author>
        <activity:object-type>http://activitystrea.ms/schema/1.0/person</activity:object-type>
        <uri>https://aaronparecki.com/</uri>
        <name>Aaron Parecki</name>
    </author>
    <link rel="alternate" href="https://aaronparecki.com/articles" type="text/html" />
    <link rel="alternate" href="https://aaronparecki.com/" type="text/html" />
    <link rel="avatar" href="https://aaronparecki.com/images/profile.jpg" />
    <link rel="self" href="https://granary.io/url?input=html&amp;output=atom&amp;url=https%3A%2F%2Faaronparecki.com%2Farticles" type="application/atom+xml" />
    <entry>
        <author>
            <activity:object-type>http://activitystrea.ms/schema/1.0/person</activity:object-type>
            <uri>https://aaronparecki.com/</uri>
            <name>Aaron Parecki</name>
        </author>
        <activity:object-type>http://activitystrea.ms/schema/1.0/article</activity:object-type>
        <id>https://aaronparecki.com/2019/02/25/9/emoji</id>
        <title>Emoji Avatars for My Website</title>
        <content type="xhtml">
            <div xmlns="http://www.w3.org/1999/xhtml"><p>My previous avatar was almost 3 years old, and I was getting tired of it. I decided to replace my avatar on my website for my IndieWebCamp Austin hack day project. But if you know me, you know I can't do anything the easy way. For my hack day project I made the avatar on each post in my website change depending on the emoji I use in the post!</p><p>I've had <a href="https://aaronparecki.com/emoji">a page on my website</a> for a while now that shows all the emoji I use in my posts. It's sorted by most frequently used emoji.</p></div>
        </content>
        <link rel="alternate" type="text/html" href="https://aaronparecki.com/2019/02/25/9/emoji" />
        <link rel="ostatus:conversation" href="https://aaronparecki.com/2019/02/25/9/emoji" />
        <activity:verb>http://activitystrea.ms/schema/1.0/post</activity:verb>
        <published>2019-02-25T09:58:24-06:00</published>
        <updated>2019-02-25T09:58:24-06:00</updated>
        <georss:point>30.269234 -97.735182</georss:point>
        <link rel="self" type="application/atom+xml" href="https://aaronparecki.com/2019/02/25/9/emoji" />
    </entry>
  </feed>
  `

	feed, err := parseAtom([]byte(doc))
	if err != nil {
		t.Error("Should not error")
	}

	if feed.Items[0].Content != `<div xmlns="http://www.w3.org/1999/xhtml"><p>My previous avatar was almost 3 years old, and I was getting tired of it. I decided to replace my avatar on my website for my IndieWebCamp Austin hack day project. But if you know me, you know I can't do anything the easy way. For my hack day project I made the avatar on each post in my website change depending on the emoji I use in the post!</p><p>I've had <a href="https://aaronparecki.com/emoji">a page on my website</a> for a while now that shows all the emoji I use in my posts. It's sorted by most frequently used emoji.</p></div>` {
		t.Error("Incorrect content found", feed.Items[0].Content)
	}
}

func TestParseContentWithCDATA(t *testing.T) {
	doc := `
  <?xml version="1.0" encoding="UTF-8"?>
<feed xml:lang="en-US" xmlns="http://www.w3.org/2005/Atom" xmlns:activity="http://activitystrea.ms/spec/1.0/" xmlns:georss="http://www.georss.org/georss" xmlns:ostatus="http://ostatus.org/schema/1.0" xmlns:thr="http://purl.org/syndication/thread/1.0" xml:base="https://aaronparecki.com/">
    <generator uri="https://granary.io/">granary</generator>
    <id>https://aaronparecki.com/articles</id>
    <title>User feed for Aaron Parecki</title>
    <logo>https://aaronparecki.com/images/profile.jpg</logo>
    <updated>2019-02-25T09:58:24-06:00</updated>
    <author>
        <activity:object-type>http://activitystrea.ms/schema/1.0/person</activity:object-type>
        <uri>https://aaronparecki.com/</uri>
        <name>Aaron Parecki</name>
    </author>
    <link rel="alternate" href="https://aaronparecki.com/articles" type="text/html" />
    <link rel="alternate" href="https://aaronparecki.com/" type="text/html" />
    <link rel="avatar" href="https://aaronparecki.com/images/profile.jpg" />
    <link rel="self" href="https://granary.io/url?input=html&amp;output=atom&amp;url=https%3A%2F%2Faaronparecki.com%2Farticles" type="application/atom+xml" />
    <entry>
        <author>
            <activity:object-type>http://activitystrea.ms/schema/1.0/person</activity:object-type>
            <uri>https://aaronparecki.com/</uri>
            <name>Aaron Parecki</name>
        </author>
        <activity:object-type>http://activitystrea.ms/schema/1.0/article</activity:object-type>
        <id>https://aaronparecki.com/2019/02/25/9/emoji</id>
        <title>Emoji Avatars for My Website</title>
        <content type="xhtml"><![CDATA[
            <div xmlns="http://www.w3.org/1999/xhtml"><p>My previous avatar was almost 3 years old, and I was getting tired of it. I decided to replace my avatar on my website for my IndieWebCamp Austin hack day project. But if you know me, you know I can't do anything the easy way. For my hack day project I made the avatar on each post in my website change depending on the emoji I use in the post!</p><p>I've had <a href="https://aaronparecki.com/emoji">a page on my website</a> for a while now that shows all the emoji I use in my posts. It's sorted by most frequently used emoji.</p></div>
        ]]></content>
        <link rel="alternate" type="text/html" href="https://aaronparecki.com/2019/02/25/9/emoji" />
        <link rel="ostatus:conversation" href="https://aaronparecki.com/2019/02/25/9/emoji" />
        <activity:verb>http://activitystrea.ms/schema/1.0/post</activity:verb>
        <published>2019-02-25T09:58:24-06:00</published>
        <updated>2019-02-25T09:58:24-06:00</updated>
        <georss:point>30.269234 -97.735182</georss:point>
        <link rel="self" type="application/atom+xml" href="https://aaronparecki.com/2019/02/25/9/emoji" />
    </entry>
  </feed>
  `

	feed, err := parseAtom([]byte(doc))
	if err != nil {
		t.Error("Should not error")
	}

	if feed.Items[0].Content != `<div xmlns="http://www.w3.org/1999/xhtml"><p>My previous avatar was almost 3 years old, and I was getting tired of it. I decided to replace my avatar on my website for my IndieWebCamp Austin hack day project. But if you know me, you know I can't do anything the easy way. For my hack day project I made the avatar on each post in my website change depending on the emoji I use in the post!</p><p>I've had <a href="https://aaronparecki.com/emoji">a page on my website</a> for a while now that shows all the emoji I use in my posts. It's sorted by most frequently used emoji.</p></div>` {
		t.Error("Incorrect content found", feed.Items[0].Content)
	}
}

func TestParseContentWithoutContent(t *testing.T) {
	doc := `
  <?xml version="1.0" encoding="UTF-8"?>
<feed xml:lang="en-US" xmlns="http://www.w3.org/2005/Atom" xmlns:activity="http://activitystrea.ms/spec/1.0/" xmlns:georss="http://www.georss.org/georss" xmlns:ostatus="http://ostatus.org/schema/1.0" xmlns:thr="http://purl.org/syndication/thread/1.0" xml:base="https://aaronparecki.com/">
    <generator uri="https://granary.io/">granary</generator>
    <id>https://aaronparecki.com/articles</id>
    <title>User feed for Aaron Parecki</title>
    <logo>https://aaronparecki.com/images/profile.jpg</logo>
    <updated>2019-02-25T09:58:24-06:00</updated>
    <author>
        <activity:object-type>http://activitystrea.ms/schema/1.0/person</activity:object-type>
        <uri>https://aaronparecki.com/</uri>
        <name>Aaron Parecki</name>
    </author>
    <link rel="alternate" href="https://aaronparecki.com/articles" type="text/html" />
    <link rel="alternate" href="https://aaronparecki.com/" type="text/html" />
    <link rel="avatar" href="https://aaronparecki.com/images/profile.jpg" />
    <link rel="self" href="https://granary.io/url?input=html&amp;output=atom&amp;url=https%3A%2F%2Faaronparecki.com%2Farticles" type="application/atom+xml" />
    <entry>
        <author>
            <activity:object-type>http://activitystrea.ms/schema/1.0/person</activity:object-type>
            <uri>https://aaronparecki.com/</uri>
            <name>Aaron Parecki</name>
        </author>
        <activity:object-type>http://activitystrea.ms/schema/1.0/article</activity:object-type>
        <id>https://aaronparecki.com/2019/02/25/9/emoji</id>
        <title>Emoji Avatars for My Website</title>
        <link rel="alternate" type="text/html" href="https://aaronparecki.com/2019/02/25/9/emoji" />
        <link rel="ostatus:conversation" href="https://aaronparecki.com/2019/02/25/9/emoji" />
        <activity:verb>http://activitystrea.ms/schema/1.0/post</activity:verb>
        <published>2019-02-25T09:58:24-06:00</published>
        <updated>2019-02-25T09:58:24-06:00</updated>
        <georss:point>30.269234 -97.735182</georss:point>
        <link rel="self" type="application/atom+xml" href="https://aaronparecki.com/2019/02/25/9/emoji" />
    </entry>
  </feed>
  `

	feed, err := parseAtom([]byte(doc))
	if err != nil {
		t.Error("Should not error")
	}

	if feed.Items[0].Content != `` {
		t.Error("Incorrect content found", feed.Items[0].Content)
	}
}
