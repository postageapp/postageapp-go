 [PostageApp](http://postageapp.com) for Go
===================================================

This is a package for Go that allows you to send emails with PostageApp service.
Personalized transactional email sending can be offloaded to PostageApp via the JSON based API.

### [API Documentation](http://help.postageapp.com/faqs/api) &bull; [PostageApp FAQs](http://help.postageapp.com/faqs) &bull; [PostageApp Help Portal](http://help.postageapp.com)

Installation
------------
In your $GOPATH/src/ directory type
<pre><code>git clone https://github.com/postageapp/postageapp-go postageapp
cd postageapp
go build
go install
</pre></code>

## Obtaining an API key

Visit [postageapp.com/register](https://secure.postageapp.com/register) and sign-up for an account. Create one or more projects
in your account each project gets its own API key. Click through to the project page and find the API key in the right-hand column.

## Sending an email

The following is a the absolute minimum required to send an email.

    cl := new(Client)
    cl.ApiKey = "YOUR_API_KEY"
    message := new(Message)
    message.Uid = uuid.Rand().Hex()

    recipient := new(Recipient)
    recipient.Email = "Alan Smithee <alan.smithee@gmail.com>"

    message.Subject = "Thank you for your order"
    message.From = "Acme Widgets <widgets@acme.com>"
    message.Recipients = append(message.Recipients, recipient)
    message.RecipientOverride = "YOUR_EMAIL_ADDRESS_HERE_DURING_DEVELOPMENT"
    message.Text = "Your order has been processed and will ship shortly."
    message.Html = <p>Your order has been processed and will ship shortly.</p>"

    response, _ := cl.SendMessage(message)

Setting the `RecipientOverride` property allows you to safely redirect all outgoing email to your own address while in development mode.

## Passing variables to templates

The real power of PostageApp kicks in when you start using templates. Templates can be configured in your PostageApp project dashboard. 
They can inherit from other templates, contain both text and html representations, provide placeholders for variables, headers and more. 
Your app doesn't need to concern itself with rendering html emails and you can update your email content without re-deploying your app. 

Once you have created a template that you want to use, specify its unique `slug` in the Template property as in the example below.

    message := new(Message)
    message.Uid = uuid.Rand().Hex()

    recipient := new(Recipient)
    recipient.Email = "Alan Smithee <alan.smithee@gmail.com>"
    message.From = "Acme Widgets <widgets@acme.com>"
    message.Template = "YOUR_TEMPLATE_SLUG" 
    message.Variables = make(map[string]string)
    message.Variables["first_name"] = "Alan"
    message.Variables["last_name"] = "Smithee"
    message.Variables["order_id"] = "555"
    message.Recipients = append(message.Recipients, recipient)

## Multiple recipients

Emails aren't restricted to just one recipient. Instead of setting the `Recipient` property, set the `Recipients` property
to a list of `Recipient` objects, each with its own set of variables.

    recipient := new(Recipient)
    recipient.Email = "Alan Smithee <alan.smithee@gmail.com>"
    recipient.Variables = make(map[string]string)
    recipient.Variables["first_name"] = "Alan"
    recipient.Variables["last_name"] = "Smithee"
    recipient.Variables["order_id"] = "555"

    recipient2 := new(Recipient)
    recipient2.Email = "Rick James <rick.james@gmail.com>"
    recipient2.Variables = make(map[string]string)
    recipient2.Variables["first_name"] = "Rick"
    recipient2.Variables["last_name"] = "James"
    recipient2.Variables["order_id"] = "556"
    message.Recipients = append(message.Recipients, recipient, recipient2)

## Attaching files

In addition to attaching files to templates in the PostageApp project dashboard, they can be attached by your app at runtime.
Simply add an `Attachment` to the `Attachments` array, providing ContentBytes, FileName and ContentType for each file attached.

    attachment := new(Attachment)
    attachment.FileName = "invoice.pdf"
    attachment.ContentType = "application/pdf"
    attachment.ContentBytes = fileBytes
    message.Attachments = append(message.Attachments, attachment)

## Adding custom headers

The `From`, `Subject` and `ReplyTo` properties are shortcuts for the following syntax.

    message.Headers = make(map[string]string)
    message.Headers["From"] = "Acme Widgets <widgets@acme.com>"
    message.Headers["Subject"] = "Your order has shipped!"
    message.Headers["Reply-To"] = "Acme Support <support@acme.com>"
    
You are free to add any necessary email headers using this method.
