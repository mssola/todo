/*
 * Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
 * This file is licensed under the MIT license.
 * See the LICENSE file.
 */


function cancelTextArea()
{
    var body = $('#contents .body');
    body.find('.contents-edit').hide();
    body.find('.contents-body').show();
}

jQuery(function() {
    $('#list button').click(function() {
        $(this).hide();

        var form = $(this).closest('.create').find('form');
        form.show();
        form.find('.tiny-text').focus();
    });

    $('#contents .cancel-btn').click(function(e) {
        e.preventDefault();
        cancelTextArea();
    });

    $('#contents .edit').click(function(e) {
        e.preventDefault();

        var body = $('#contents .body');
        body.find('.contents-body').hide();
        body.find('.contents-edit').show();
        $('textarea').focus();
    });

    $('#contents .rename').click(function(e) {
        e.preventDefault();

        var f = $('#contents .header').find('.rename-form');
        f.show();
        f.find('.tiny-text').focus();
    });

    $('#contents .delete').click(function(e) {
        e.preventDefault();
        $(this).hide();
        $('#contents .confirmation').css('display', 'inline');
    });

    $('#contents .yes').click(function(e) {
        e.preventDefault();
        $('#contents .delete-form').submit();
    });

    $('#contents .no').click(function(e) {
        e.preventDefault();
        $('#contents .confirmation').hide();
        $('#contents .delete').css('display', 'inline');
    });

    $('#contents textarea').focusout(cancelTextArea);
});

