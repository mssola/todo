
# How to build jQuery

I'm using a custom build of jQuery in this project. This is because I need
just a few things of jQuery, but not the whole pack. Let's see how we can
build jQuery.

First of all, we need to download it. We do this with git and select the
version 2.1.1.

    cd /tmp
    git clone https://github.com/jquery/jquery.git && cd jquery
    git checkout tags/2.1.1

After this, we need to install grunt. Inside the jquery directory we have
to perform the following commands.

    sudo npm install -g grunt-cli
    npm install

And finally, we use grunt to build jQuery. Note that in the following command
we use the `custom` option. With this option we specify the modules that we do
**not** want. Without further due, this is the command to be executed:

    grunt custom -ajax,-css,-deprecated,-dimensions,-effects,-event/alias,-offset,-wrap,-exports/amd

The result will be available inside the `dist` directory.
