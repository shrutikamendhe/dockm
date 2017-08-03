angular.module('dockm.services')
.factory('ModalService', [function ModalServiceFactory() {
  'use strict';
  var service = {};

  var applyBoxCSS = function(box) {
    box.css({
      'top': '50%',
      'margin-top': function () {
        return -(box.height() / 2);
      }
    });
  };

  var confirmButtons = function(options) {
    var buttons = {
      confirm: {
        label: options.buttons.confirm.label,
        className: options.buttons.confirm.className
      },
      cancel: {
        label: options.buttons.cancel && options.buttons.cancel.label ? options.buttons.cancel.label : 'Cancel'
      }
    };
    return buttons;
  };

  service.confirm = function(options){
    var box = bootbox.confirm({
      title: options.title,
      message: options.message,
      buttons: confirmButtons(options),
      callback: options.callback
    });
    applyBoxCSS(box);
  };

  service.prompt = function(options){
    var box = bootbox.prompt({
      title: options.title,
      inputType: options.inputType,
      inputOptions: options.inputOptions,
      buttons: confirmButtons(options),
      callback: options.callback
    });
    applyBoxCSS(box);
  };

  service.confirmAccessControlUpdate = function(callback, msg) {
    service.confirm({
      title: 'Are you sure ?',
      message: 'Changing the ownership of this resource will potentially restrict its management to some users.',
      buttons: {
        confirm: {
          label: 'Change ownership',
          className: 'btn-primary'
        }
      },
      callback: callback
    });
  };

  service.confirmImageForceRemoval = function(callback) {
    service.confirm({
      title: 'Are you sure?',
      message: 'Forcing the removal of the image will remove the image even if it has multiple tags or if it is used by stopped containers.',
      buttons: {
        confirm: {
          label: 'Remove the image',
          className: 'btn-danger'
        }
      },
      callback: callback
    });
  };

  service.confirmDeletion = function(message, callback) {
    service.confirm({
      title: 'Are you sure ?',
      message: message,
      buttons: {
        confirm: {
          label: 'Remove',
          className: 'btn-danger'
        }
      },
      callback: callback
    });
  };

  service.confirmContainerDeletion = function(title, callback) {
    service.prompt({
      title: title,
      inputType: 'checkbox',
      inputOptions: [
        {
          text: 'Automatically remove non-persistent volumes<i></i>',
          value: '1'
        }
      ],
      buttons: {
        confirm: {
          label: 'Remove',
          className: 'btn-danger'
        }
      },
      callback: callback
    });
  };

  return service;
}]);
