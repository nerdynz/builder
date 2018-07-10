;(function (window, $) {
  'use strict'
  $.ImageUpload = function ImageUpload (options, element) {
    this.element = element
    this.options = $.extend(
      true,
      {},
      $.ImageUpload.defaults,
      options,
      $(this.element).data()
    )
    this.input = $(this.element).find('input[type=file]')

    var _self = this
    this.interval = null
    this.drag = false
    this.widthPercentage = 100

    // buttons
    this.button = {}
    this.button.edit =
      '<div class="control"><div class="button is-small is-info button-edit" title="' +
      (this.options.editTitle || 'Edit') +
      '"><i class="fa fa-pencil"></i></div></div>'
    this.button.saving =
      '<div class="control"><div class="button is-small is-warning saving">' +
      (this.options.saveLabel || 'Saving...') +
      ' <i class="fa fa-time"></i></div></div>'
    this.button.zoomin =
      '<div class="control"><div class="button is-small button-zoom-in" title="' +
      (this.options.zoominTitle || 'Zoom in') +
      '"><i class="fa fa-search-plus"></i></div></div>'
    this.button.zoomout =
      '<div class="control"><div class="button is-small button-zoom-out" title="' +
      (this.options.zoomoutTitle || 'Zoom out') +
      '"><i class="fa fa-search-minus"></i></div></div>'
    this.button.zoomreset =
      '<div class="control"><div class="button is-small button-zoom-reset" title="' +
      (this.options.zoomresetTitle || 'Reset') +
      '"><i class="fa fa-refresh"></i></div></div>'
    // this.button.rotatecw = '<div class="control"><div class="button is-small button-rotate-cw" title="' + (this.options.cwTitle || 'Rotate clockwise') + '"><i class="fa fa-share"></i></div></div>';
    // this.button.rotateccw= '<div class="control"><div class="button is-small button-rotate-ccw" title="' + (this.options.ccwTitle || 'Rotate counter clockwise') + '"><i class="fa fa-share icon-flipped"></i></div></div>';
    this.button.cancel =
      '<div class="control"><div class="button is-small is-danger button-cancel" title="' +
      (this.options.cancelTitle || 'Cancel') +
      '"><i class="fa fa-remove"></i></div></div>'
    this.button.done =
      '<div class="control"><div class="button is-small is-success button-ok" title="' +
      (this.options.okTitle || 'Ok') +
      '"><i class="fa fa-check"></i></div></div>'
    this.button.del =
      '<div class="control"><div class="button is-small is-danger button-del" title="' +
      (this.options.delTitle || 'Delete') +
      '"><i class="fa fa-trash"></i></div></div>'

    this.button.download =
      '<a class="button is-small button-warning download"><i class="fa fa-download"></i> ' +
      (this.options.downloadLabel || 'Download') +
      '</a>'

    this.button.fullscreen = '<div class="control"><div class="button is-small button-zoom-reset" title="' +
      (this.options.fullscreenTitle || 'Fullscreen') +
      '"><i class="fa fa-arrows-alt"></i></div></div>'
    // this.imageExtensions = {png: 'png', bmp: 'bmp', gif: 'gif', jpg: ['jpg','jpeg'], tiff: 'tiff'};
    this.imageMimes = {
      png: 'image/png',
      bmp: 'image/bmp',
      gif: 'image/gif',
      jpg: 'image/jpeg',
      jpeg: 'image/jpeg',
      tiff: 'image/tiff'
    }

    _self._init()
  }

  $.ImageUpload.defaults = {
    width: null,
    height: null,
    image: null,
    ghost: true,
    originalsize: true,
    url: false,
    removeurl: null,
    data: {},
    canvas: true,
    canvasImageOnly: false,
    ajax: true,
    resize: false,
    dimensionsonly: false,
    editstart: false,
    saveOriginal: false,
    save: true,
    download: false,

    smaller: false,
    smallerWidth: false,
    smallerHeight: false,
    background: null,

    onAfterZoomImage: null,
    onAfterInitImage: null,
    onAfterMoveImage: null,
    onAfterProcessImage: null,
    onAfterResetImage: null,
    onAfterCancel: null,
    onAfterRemoveImage: null,
    onAfterSelectImage: null
  }

  $.ImageUpload.prototype = {
    _init: function () {
      var _self = this
      var options = _self.options
      var element = _self.element
      var input = _self.input

      if (empty($(element))) {
        return false
      } else {
        $(element).children().css({ position: 'absolute' })
      }

      // the engine of this script
      if (!(window.FormData && 'upload' in $.ajaxSettings.xhr())) {
        $(element)
          .empty()
          .attr('class', '')
          .addClass('alert alert-danger')
          .html(
            'HTML5 Upload Image: Sadly.. this browser does not support the plugin, update your browser today!'
          )
        return false
      }

      if (
        !empty(options.width) &&
        empty(options.height) &&
        $(element).innerHeight() <= 0
      ) {
        $(element)
          .empty()
          .attr('class', '')
          .addClass('alert alert-danger')
          .html(
            'HTML5 Upload Image: Image height is not set and can not be calculated...'
          )
        return false
      }
      if (
        !empty(options.height) &&
        empty(options.width) &&
        $(element).innerWidth() <= 0
      ) {
        $(element)
          .empty()
          .attr('class', '')
          .addClass('alert alert-danger')
          .html(
            'HTML5 Upload Image: Image width is not set and can not be calculated...'
          )
        return false
      }
      /* if (!empty(options.height) && !empty(options.width) && !empty($(element).innerHeight() && !empty($(element).innerWidth()))) {
       //all sizes are filled in
       console.log(options.width)
       console.log(options.height)
       console.log(options.width / options.height)

       console.log($(element).outerWidth())
       console.log($(element).outerHeight())

       console.log($(element).outerWidth() / $(element).outerHeight())

       if ((options.width / options.height) !== ($(element).outerWidth() / $(element).outerHeight())) {
         $(element).empty().attr('class','').addClass('alert alert-danger').html('HTML5 Upload Image: The ratio of all sizes (CSS and image) are not the same...');
         return false;
       }
      } */

      // copy styles
      $(element).data('style', $(element).attr('style'))
      $(element).data('class', $(element).attr('class'))

      /// ///////////
      // check sizes
      var width, height, imageWidth, imageHeight

      options.width = imageWidth = options.width || $(element).outerWidth()
      options.height = imageHeight = options.height || $(element).outerHeight()

      if ($(element).innerWidth() > 0) {
        width = $(element).outerWidth()
      } else if ($(element).innerHeight() > 0) {
        width = null
      } else if (!empty(options.width)) {
        width = options.width
      }

      if ($(element).innerHeight() > 0) {
        height = $(element).outerHeight()
      } else if ($(element).innerWidth() > 0) {
        height = null
      } else if (!empty(options.height)) {
        height = options.height
      }

      height = height || width / (imageWidth / imageHeight)
      width = width || height / (imageHeight / imageWidth)

      /* is small window, add class small */
      if (width < 240) {
        $(element).addClass('smalltools smalltext')
      }

      $(element).css({ height: height, width: width })
      _self.widthPercentage =
        $(element).outerWidth() / $(element).offsetParent().width() * 100

      if (options.meta) {
        _self.resize(options.meta.containerWidth, options.meta.containerHeight)
      }
      if (options.resize === true) {
        $(window).resize(function () {
          _self.resize()
        })
      }
      _self._bind()

      if (options.required || $(input).prop('required')) {
        _self.options.required = true
        $(input).prop('required', true)
      }

      if (!options.ajax) {
        _self._formValidation()
      }
      if (options.meta && options.meta.originalName) {
        _self.readImage(
          options.meta.originalName,
          options.meta.originalName,
          options.meta.originalName,
          _self.imageMimes[options.meta.originalName.split('.').pop()],
          options.meta,
          function () {
            _self.imageCrop(true, function () {
              _self.resize()
              if (_self.options.onAfterInitImage) { _self.options.onAfterInitImage.call(_self, _self) }
            })
          }
        ) // $(img).attr('src'),)
      } else if (!empty(options.image) && options.editstart === false) {
        $(element)
          .data('name', options.image)
          .append($('<img />').addClass('preview').attr('src', options.image))

        var tools = $('<div class="preview tools field has-addons"></div>')
        var del = $('' + this.button.del + '')
        /* $(edit).unbind('click').click(function(e) {
         e.preventDefault();
          $(element).children().show();
          $(element).find('.final').remove();
          $(input).data('valid',false);
        }) */
        $(del).unbind('click').click(function (e) {
          e.preventDefault()
          _self.reset()
        })

        // var edit = $('' + this.button.edit + '')
        // $(edit).unbind('click').click(function (e) {
        //   e.preventDefault()
        //   _self.reset()
        // })

        // if (options.buttonEdit !== false) {
        //   $(tools).append(edit)
        // }

        if (options.buttonDel !== false) {
          $(tools).append(del)
        }

        $(element).append(tools)
        if (_self.options.onAfterInitImage) { _self.options.onAfterInitImage.call(_self, _self) }
      } else if (!empty(options.image)) {
        _self.readImage(
          options.image,
          options.image,
          options.image,
          _self.imageMimes[options.image.split('.').pop()],
          null,
          function () {
            if (_self.options.onAfterInitImage) { _self.options.onAfterInitImage.call(_self, _self) }
          }
        ) // $(img).attr('src'),)
      } else {
        if (_self.options.onAfterInitImage) { _self.options.onAfterInitImage.call(_self, _self) }
      }
    },
    _bind: function () {
      var _self = this
      var element = _self.element
      var input = _self.input

      // bind the events
      $(element).unbind('dragover drop mouseout').on({
        dragover: function (event) {
          _self.handleDrag(event)
        },
        drop: function (event) {
          _self.handleFile(event, $(this))
        },
        mouseout: function () {
          $(document)
            .unbind('mouseup touchend')
            .on('mouseup touchend', function (e) {
              e.preventDefault()
              _self.imageUnhandle() //
            })
        }
      })

      $(input).unbind('change').change(function (event) {
        _self.drag = false
        _self.handleFile(event, $(element))
      })
    },

    handleFile: function (event, element) {
      event.stopPropagation()
      event.preventDefault()

      var _self = this
      var options = this.options
      var files = _self.drag === false
        ? event.originalEvent.target.files
        : event.originalEvent.dataTransfer.files // FileList object.
      _self.drag = false

      // _self.reset();
      if (options.removeurl !== null && !empty($(element).data('name'))) {
        $.ajax({
          type: 'POST',
          url: options.removeurl,
          data: { image: $(element).data('name') },
          success: function (response) {
            if (_self.options.onAfterRemoveImage) { _self.options.onAfterRemoveImage.call(_self, response, _self) }
          }
        })
      }

      $(element).removeClass('notAnImage').addClass('loading') // .unbind('dragover').unbind('drop');

      for (var i = 0, f; (f = files[i]); i++) {
        if (!f.type.match('image.*')) {
          $(element).addClass('notAnImage')
          continue
        }

        var reader = new window.FileReader()

        reader.onload = (function (theFile) {
          // console.log(theFile);
          return function (e) {
            $(element).find('img').remove()
            _self.readImage(
              reader.result,
              e.target.result,
              theFile.name,
              theFile.type
            )
          }
        })(f)
        reader.readAsDataURL(f)
      }
      if (_self.options.onAfterSelectImage) { _self.options.onAfterSelectImage.call(_self, null, _self) }
    },

    readImage: function (image, src, name, mimeType, serverMeta, cb) {
      var _self = this
      var options = this.options
      var element = this.element
      _self.drag = false

      var img = new window.Image()
      img.onload = function (tmp) {
        var imgElement = $('<img src="' + src + '" name="' + name + '" />')
        var width, height, useWidth, useHeight, ratio, elementRatio
        if (serverMeta) {
          width = useWidth = serverMeta.cssWidth
          height = useHeight = serverMeta.cssHeight
        } else {
          width = useWidth = img.width
          height = useHeight = img.height
        }
        ratio = width / height
        elementRatio = $(element).outerWidth() / $(element).outerHeight()
        // resize image
        if (options.originalsize === false) {
          // need to add the option is smaller = true, dont resize
          // also if the image === 100% the size of the element, dont add extra space

          useWidth = $(element).outerWidth() + 40
          useHeight = useWidth / ratio

          if (useHeight < $(element).outerHeight()) {
            useHeight = $(element).outerHeight() + 40
            useWidth = useHeight * ratio
          }
        } else if (
          useWidth < $(element).outerWidth() ||
          useHeight < $(element).outerHeight()
        ) {
          if (options.smaller === true) {
            // do nothing
          } else {
            if (ratio < elementRatio) {
              useWidth = $(element).outerWidth()
              useHeight = useWidth / ratio
            } else {
              useHeight = $(element).outerHeight()
              useWidth = useHeight * ratio
            }
          }
        }

        var left, top
        if (serverMeta) {
          left = serverMeta.cssLeft + 'px'
          top = serverMeta.cssTop + 'px'
        } else {
          left = parseFloat(($(element).outerWidth() - useWidth) / 2) // * -1;
          top = parseFloat(($(element).outerHeight() - useHeight) / 2) // * -1;
        }

        imgElement.css({
          left: left,
          top: top,
          width: useWidth,
          height: useHeight
        })
        _self.image = $(imgElement)
          .clone()
          .data({
            mime: mimeType,
            width: width,
            height: height,
            ratio: ratio,
            left: left,
            top: top,
            useWidth: useWidth,
            useHeight: useHeight
          })
          .addClass('main')
          .bind('mousedown touchstart', function (event) {
            _self.imageHandle(event)
          })
        _self.imageGhost = options.ghost
          ? $(imgElement).addClass('ghost')
          : null

        // place the images
        $(element).append(
          $('<div class="cropWrapper"></div>').append($(_self.image))
        )
        if (!empty(_self.imageGhost)) {
          $(element).append(_self.imageGhost)
        }

        // $(element).unbind('dragover').unbind('drop');
        _self._tools()

        // clean up
        $(element).removeClass('loading')
        if (cb) { cb() }
      }
      img.onerror = function () {
        if (options.onLoadFailed) {
          options.onLoadFailed()
        }
      }
      img.src = image
    },

    handleDrag: function (event) {
      var _self = this
      _self.drag = true
      event.stopPropagation()
      event.preventDefault()
      event.originalEvent.dataTransfer.dropEffect = 'copy'
    },
    imageHandle: function (e) {
      e.preventDefault() // disable selection
      var event = e.originalEvent.touches || e.originalEvent.changedTouches
        ? e.originalEvent.touches[0] || e.originalEvent.changedTouches[0]
        : e

      var _self = this
      var element = this.element
      var image = this.image
      var options = this.options

      var height = image.outerHeight()
      var width = image.outerWidth()
      var cursorY = image.offset().top + height - event.pageY
      var cursorX = image.offset().left + width - event.pageX

      $(document).on({
        'dragstart mousemove touchmove': function (e) {
          $('body').css({ cursor: 'move' })

          e.stopImmediatePropagation()
          e.preventDefault()

          var event = e.originalEvent.touches || e.originalEvent.changedTouches
            ? e.originalEvent.touches[0] || e.originalEvent.changedTouches[0]
            : e

          var imgTop = event.pageY + cursorY - height
          var imgLeft = event.pageX + cursorX - width
          var hasBorder = $(element).outerWidth() !== $(element).innerWidth()

          if (options.smaller === false) {
            if (parseInt(imgTop - $(element).offset().top) > 0) {
              imgTop = $(element).offset().top + (hasBorder ? 1 : 0)
            } else if (
              imgTop + height <
              $(element).offset().top + $(element).outerHeight()
            ) {
              imgTop =
                $(element).offset().top +
                $(element).outerHeight() -
                height +
                (hasBorder ? 1 : 0)
            }

            if (parseInt(imgLeft - $(element).offset().left) > 0) {
              imgLeft = $(element).offset().left + (hasBorder ? 1 : 0)
            } else if (
              imgLeft + width <
              $(element).offset().left + $(element).outerWidth()
            ) {
              imgLeft =
                $(element).offset().left +
                $(element).outerWidth() -
                width +
                (hasBorder ? 1 : 0)
            }
          }

          image.offset({
            top: imgTop,
            left: imgLeft
          })
          _self._ghost()
        },
        'dragend mouseup touchend': function () {
          _self.imageUnhandle()
        }
      })
    },
    imageUnhandle: function () {
      var _self = this
      // var image = _self.image

      $(document).unbind('mousemove touchmove')
      $('body').css({ cursor: '' })
      if (_self.options.onAfterMoveImage) { _self.options.onAfterMoveImage.call(_self, _self) }
    },
    imageZoom: function (x) {
      var _self = this
      var element = _self.element
      var image = _self.image
      var options = _self.options

      if (empty(image)) {
        _self._clearTimers()
        return false
      }

      var ratio = image.data('ratio')
      var newWidth = image.outerWidth() + x
      var newHeight = newWidth / ratio

      if (options.smaller === false) {
        // smaller then element?
        if (newWidth < $(element).outerWidth()) {
          newWidth = $(element).outerWidth()
          newHeight = newWidth / ratio
        }
        if (newHeight < $(element).outerHeight()) {
          newHeight = $(element).outerHeight()
          newWidth = newHeight * ratio
        }
      } else {
        if (newWidth < 1 || newHeight < 1) {
          if (newWidth > newHeight) {
            newWidth = 1
            newHeight = newWidth / ratio
          } else {
            newHeight = 1
            newWidth = newHeight * ratio
          }
        }
      }

      var newTop =
        image.css('top').replace(/[^-\d.]/g, '') -
        (newHeight - image.outerHeight()) / 2
      var newLeft =
        image.css('left').replace(/[^-\d.]/g, '') -
        (newWidth - image.outerWidth()) / 2

      if (options.smaller === false) {
        if ($(element).offset().left - newLeft < $(element).offset().left) {
          newLeft = 0
        } else if (
          $(element).outerWidth() > newLeft + $(image).outerWidth() &&
          x <= 0
        ) {
          newLeft = $(element).outerWidth() - newWidth
        }

        if ($(element).offset().top - newTop < $(element).offset().top) {
          newTop = 0
        } else if (
          $(element).outerHeight() > newTop + $(image).outerHeight() &&
          x <= 0
        ) {
          newTop = $(element).outerHeight() - newHeight
        }
      }
      image.css({
        width: newWidth,
        height: newHeight,
        top: newTop,
        left: newLeft
      })
      _self._ghost()
    },
    imageCrop: function (skipSave, cb) {
      var _self = this
      var element = _self.element
      var image = _self.image
      var input = _self.input
      var options = _self.options

      var factor = options.width !== $(element).outerWidth()
        ? options.width / $(element).outerWidth()
        : 1

      var finalWidth,
        finalHeight,
        finalTop,
        finalLeft,
        imageWidth,
        imageHeight,
        imageOriginalWidth,
        imageOriginalHeight

      finalWidth = options.width
      finalHeight = options.height

      finalTop = parseInt(parseInt($(image).css('top')) * factor) + 1
      finalLeft = parseInt(parseInt($(image).css('left')) * factor) + 1
      imageWidth = parseInt($(image).width() * factor)
      imageHeight = parseInt($(image).height() * factor)
      imageOriginalWidth = parseInt($(image).data('width'))
      imageOriginalHeight = parseInt($(image).data('height'))

      finalTop = finalTop || 0
      finalLeft = finalLeft || 0
      var obj = {
        name: $(image).attr('name'),
        imageOriginalWidth: imageOriginalWidth,
        imageOriginalHeight: imageOriginalHeight,
        imageWidth: imageWidth,
        imageHeight: imageHeight,
        width: finalWidth,
        height: finalHeight,
        left: finalLeft,
        top: finalTop,
        cssLeft: parseFloat($(image).css('left')),
        cssTop: parseFloat($(image).css('top')),
        cssWidth: parseFloat($(image).css('width')),
        cssHeight: parseFloat($(image).css('height')),
        containerWidth: $(element).outerWidth(),
        containerHeight: $(element).outerHeight()
      }

      $(element).find('.tools').children().toggle()
      if (!skipSave) {
        $(element).find('.tools').append($(_self.button.saving))
      }

      if (options.canvas === true) {
        var canvas = $(
          '<canvas class="final" id="canvas_' +
            $(input).attr('name') +
            '" width="' +
            finalWidth +
            '" height="' +
            finalHeight +
            '" style="position:absolute; top: -1px; bottom: 0; left:  -1px; right: 0; z-index:100; width: calc(100% + 2px); height: calc(100% + 2px);"></canvas>'
        )

        $(element).append(canvas)

        var canvasContext = $(canvas)[0].getContext('2d')
        var imageObj = new window.Image()

        // canvasContext.fillStyle = "rgba(255, 255, 255, 0)";
        // canvasContext.clearRect(0,0,finalWidth,finalHeight);

        imageObj.onload = function () {
          var canvasTmp = $(
            '<canvas width="' +
              imageWidth +
              '" height="' +
              imageHeight +
              '"></canvas>'
          )
          var canvasTmpContext = $(canvasTmp)[0].getContext('2d')

          // canvasTmpContext.fillStyle = "rgba(255, 255, 255, 0)";
          // canvasTmpContext.clearRect(0,0,imageWidth,imageHeight);
          canvasTmpContext.drawImage(imageObj, 0, 0, imageWidth, imageHeight)
          // $(element).append(canvasTmp);
          var tmpObj = new window.Image()
          tmpObj.onload = function () {
            if (options.canvasImageOnly === true) {
              var _imageWidth = imageWidth
              var _imageHeight = imageHeight
              if (finalLeft < 0) {
                _imageWidth += finalLeft
              } else if (finalLeft + imageWidth > finalWidth) {
                _imageWidth = finalWidth - finalLeft
              }
              if (finalTop < 0) {
                _imageHeight += finalTop
              } else if (finalTop + imageHeight > finalHeight) {
                _imageHeight = finalHeight - finalTop
              }

              var canvasImageOnly = $(
                '<canvas width="' +
                  _imageWidth +
                  '" height="' +
                  _imageHeight +
                  '"></canvas>'
              )
              var canvasImageOnlyContext = $(canvasImageOnly)[0].getContext(
                '2d'
              )
              canvasImageOnlyContext.drawImage(
                tmpObj,
                finalLeft < 0 ? finalLeft : 0,
                finalTop < 0 ? finalTop : 0,
                imageWidth,
                imageHeight
              )
            }

            if (imageWidth < finalWidth || imageHeight < finalHeight) {
              canvasContext.drawImage(
                tmpObj,
                finalLeft,
                finalTop,
                imageWidth,
                imageHeight
              )
            } else {
              canvasContext.drawImage(
                tmpObj,
                (finalLeft * -1),
                finalTop * -1,
                finalWidth,
                finalHeight,
                0,
                0,
                finalWidth,
                finalHeight
              )
            }

            var dataUrl = options.canvasImageOnly === true
              ? $(canvasImageOnly)[0].toDataURL(image.data('mime'))
              : $(canvas)[0].toDataURL(image.data('mime'))

            if (skipSave) {
              $(element).find('.tools .saving').remove()
              $(element).find('.tools').children().toggle()
              _self.imageFinal(cb)
            } else if (options.save === false) {
              $(element).find('.tools .saving').remove()
              $(element).find('.tools').children().toggle()

              if (_self.options.onSave) {
                _self.options.onSave.call(
                  _self,
                  $.extend(obj, options.data, { data: dataUrl })
                )
              }
              _self.imageFinal(cb)
            } else if (options.ajax === true) {
              _self._ajax($.extend({ data: dataUrl }, obj))
            } else {
              var json = JSON.stringify($.extend({ data: dataUrl }, obj))
              $(input).after(
                $(
                  '<input type="text" name="' +
                    $(input).attr('name') +
                    '_values" class="final" />'
                ).val(json)
              )

              $(input).data('required', $(input).prop('required'))
              $(input).prop('required', false)
              $(input).wrap('<form>').parent('form').trigger('reset')
              $(input).unwrap()

              $(element).find('.tools .saving').remove()
              $(element).find('.tools').children().toggle()

              _self.imageFinal(cb)
            }
          }
          tmpObj.src = $(canvasTmp)[0].toDataURL(image.data('mime'))
        }
        imageObj.src = $(image).attr('src')

        if (options.saveOriginal === true) {
          // console.log($(image).attr('src'));
          obj = $.extend({ original: $(image).attr('src') }, obj)
        }
      } else {
        if (options.ajax === true) {
          _self._ajax(
            $.extend(
              {
                data: $(image).attr('src'),
                saveOriginal: options.saveOriginal
              },
              obj
            )
          )
        } else {
          var finalImage = $(element).find('.cropWrapper').clone()
          $(finalImage).addClass('final').show().unbind().children().unbind()
          $(element).append($(finalImage))

          _self.imageFinal(cb)

          var json = JSON.stringify(obj)
          $(input).after(
            $(
              '<input type="text" name="' +
                $(input).attr('name') +
                '_values" class="final" />'
            ).val(json)
          )
        }
      }
    },
    _ajax: function (obj) {
      var _self = this
      var element = _self.element
      // var image = _self.image
      var options = _self.options

      if (options.dimensionsonly === true) {
        delete obj.data
      }

      $.ajax({
        type: 'POST',
        url: options.url,
        data: $.extend(obj, options.data),
        success: function (response) {
          if (response.status === 'success') {
            var file = response.url.split('?')
            $(element).find('.tools .saving').remove()
            $(element).find('.tools').children().toggle()
            $(element).data('name', file[0])
            $(element).data('filename', response.filename)
            if (options.canvas !== true) {
              $(element).append(
                $(
                  '<img src="' +
                    file[0] +
                    '" class="final" style="width: 100%" />'
                )
              )
            }

            _self.imageFinal()
          } else {
            $(element).find('.tools .saving').remove()
            $(element).find('.tools').children().toggle()
            $(element).append(
              $(
                '<div class="alert alert-danger">' + response.error + '</div>'
              ).css({
                bottom: '10px',
                left: '10px',
                right: '10px',
                position: 'absolute',
                zIndex: 99
              })
            )
            setTimeout(function () {
              _self.responseReset()
            }, 2000)
          }
        },
        error: function (response, status) {
          $(element).find('.tools .saving').remove()
          $(element).find('.tools').children().toggle()
          $(element).append(
            $(
              '<div class="alert alert-danger"><strong>' +
                response.status +
                '</strong> ' +
                response.statusText +
                '</div>'
            ).css({
              bottom: '10px',
              left: '10px',
              right: '10px',
              position: 'absolute',
              zIndex: 99
            })
          )
          setTimeout(function () {
            _self.responseReset()
          }, 2000)
        }
      })
    },
    imageReset: function () {
      var _self = this
      var image = _self.image
      // var element = _self.element

      $(image).css({
        width: image.data('useWidth'),
        height: image.data('useHeight'),
        top: image.data('top'),
        left: image.data('left')
      })
      _self._ghost()

      if (_self.options.onAfterResetImage) { _self.options.onAfterResetImage.call(_self, _self) }
    },
    imageFinal: function (cb) {
      var _self = this
      var element = _self.element
      var input = _self.input
      var options = _self.options

      // remove all children except final
      $(element).addClass('done')
      $(element).children().not('.final').hide()

      // create tools element
      var tools = $('<div class="tools final field has-addons">')

      // edit option after crop
      if (options.meta.cantEdit) {
        // dont do anything
      } else if (options.buttonEdit !== false) {
        $(tools).append(
          $(_self.button.edit).click(function () {
            $(element).children().show()
            $(element).find('.final').remove()
            $(input).data('valid', false)
          })
        )
      }

      // delete option after crop
      if (options.buttonDel !== false) {
        $(tools).append(
          $(_self.button.del).click(function (e) {
            _self.reset()
          })
        )
      }

      // if in fullscreen mode then close
      if ($(element).data('isFullscreen')) {
        _self.expandAndResize()
      }

      // append tools to element
      $(element).append(tools)
      $(element).unbind()

      // download
      if (options.download === true) {
        var download = $('<div class="download final"/>')
        $(download).append(
          $(_self.button.download)
            .attr('download', $(_self.image).attr('name'))
            .click(function (e) {
              $(this).attr(
                'href',
                $(element)
                  .find('canvas.final')[0]
                  .toDataURL(_self.image.data('mime'))
              )
            })
        )
        $(element).append(download)
      }

      // set input to valid for form upload
      $(input).unbind().data('valid', true)

      // custom function after process image;
      if (cb) { cb() }
      if (_self.options.onAfterProcessImage) { _self.options.onAfterProcessImage.call(_self, _self) }
    },
    responseReset: function () {
      var _self = this
      var element = _self.element

      // remove responds from ajax event
      $(element).find('.alert').remove()
    },
    reset: function () {
      var _self = this
      var element = _self.element
      var input = _self.input
      var options = _self.options
      _self.image = null

      $(element).find('.preview').remove()
      $(element)
        .removeClass('loading done')
        .children()
        .show()
        .not('input[type=file]')
        .remove()
      $(input).wrap('<form>').parent('form').trigger('reset')
      $(input).unwrap()
      $(input)
        .prop(
          'required',
          $(input).data('required') || options.required || false
        )
        .data('valid', false)
      _self._bind()

      if (options.removeurl !== null && !empty($(element).data('name'))) {
        $.ajax({
          type: 'POST',
          url: options.removeurl,
          data: { image: $(element).data('name') },
          success: function (response) {
            if (_self.options.onAfterRemoveImage) { _self.options.onAfterRemoveImage.call(_self, response, _self) }
          }
        })
      } else {
        if (_self.options.onAfterRemoveImage) { _self.options.onAfterRemoveImage.call(_self, null, _self) }
      }
      $(element).data('name', null)

      if (_self.imageGhost) {
        $(_self.imageGhost).remove()
        _self.imageGhost = null
      }

      if ($(element).data('isFullscreen')) {
        _self.expandAndResize()
      }

      if (_self.options.onAfterCancel) _self.options.onAfterCancel.call(_self)
    },
    expandAndResize: function () {
      var _self = this
      var element = $(_self.element)
      var overlay = element.prev()
      // var image = _self.image
      if (element.data('isFullscreen')) {
        element.data('isFullscreen', false)
        overlay.hide()
        element.css({
          position: 'initial',
          top: 'initial',
          left: 'initial',
          right: 'initial',
          bottom: 'initial',
          zIndex: 2000
        })
      } else {
        overlay.show()
        element.data('isFullscreen', true)
        element.css({
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          zIndex: 2000
        })
      }
      _self.resize(null, null, true)
    },
    resize: function (containerWidthOverride, containerHeightOverride, isFullscreen) {
      var _self = this
      var options = _self.options
      var element = _self.element
      var image = _self.image

      if (options.resize !== true) return false

      var containerWidth = $(element).outerWidth()
      var parentWidth = $(element).offsetParent().width()
      var parentHeight = $(element).offsetParent().height()

      var width = parentWidth * (_self.widthPercentage / 100)
      if (width > options.width && isFullscreen) {
        width = options.width
      }
      var factor = width / containerWidth
      var height = $(element).outerHeight() * factor
      if (width > options.width && isFullscreen) {
        height = options.height * factor
      }
      if (containerWidthOverride && containerHeightOverride) {
        width = containerWidthOverride
        height = containerHeightOverride
      }
      // $(element).css({ width: containerWidthOverride, height: containerHeightOverride })
      if (isFullscreen) {
        $(element).css({ width: width, height: height, left: (parentWidth - width) / 2, top: (parentHeight - height) / 2 })
      } else {
        $(element).css({ width: width, height: height })
      }

      if (width < 240) {
        if (!$(element).hasClass('smalltools smalltext')) {
          $(element).addClass('smalltools smalltext smalladded')
        }
      } else {
        if ($(element).hasClass('smalladded')) {
          $(element).removeClass('smalltools smalltext smalladded')
        }
      }

      if (!empty(image)) {
        // console.log(image.offset());
        $(image).css({
          left: $(image).css('left').replace(/[^-\d.]/g, '') * factor + 'px',
          top: $(image).css('top').replace(/[^-\d.]/g, '') * factor + 'px'
        })
        $(image).width($(image).width() * factor)
        $(image).height($(image).height() * factor)

        _self._ghost()
      }
      // console.log('resize plugin');
      return true
    },
    // rotate: function (degrees) {
    //   var _self = this
    //   var element = _self.element
    //   var image = _self.image

    //   $(image).addClass('rotate_90')
    //   var tmp = $(image).data('width')
    //   $(image).data('width', $(image).data('height'))
    //   $(image).data('height', tmp)
    // },
    reinit: function () {
      var _self = this
      // var element = _self.element
      _self.resize()
      _self._bind()
      return true
    },
    modal: function () {
      var _self = this
      var element = _self.element

      $(element).attr('style', $(element).data('style'))
      $(element).attr('class', $(element).data('class'))

      _self._init()

      return this
    },

    _ghost: function () {
      var _self = this
      var options = _self.options
      var image = _self.image
      var ghost = _self.imageGhost

      // if set to true, mirror all drag events
      // function in one place, much needed
      if (options.ghost === true && !empty(ghost)) {
        $(ghost).css({
          width: image.css('width'),
          height: image.css('height'),
          top: image.css('top'),
          left: image.css('left')
        })
      }
    },
    _tools: function () {
      var _self = this
      var element = _self.element
      var tools = $('<div class="tools field has-addons"></div>')
      var options = _self.options

      // zoomin button
      if (options.buttonZoomin !== false) {
        $(tools).append(
          $(_self.button.zoomin).on({
            'touchstart mousedown': function (e) {
              e.preventDefault()
              var imageZoom = 0.5
              var zooming = 0
              _self.interval = window.setInterval(function () {
                zooming++
                if (zooming % 150 === 0) {
                  imageZoom = imageZoom + 0.5
                }
                _self.imageZoom(imageZoom)
              }, 1)
            },
            'touchend mouseup mouseleave': function (e) {
              e.preventDefault()
              window.clearInterval(_self.interval)
              if (_self.options.onAfterZoomImage) { _self.options.onAfterZoomImage.call(_self, _self) }
            }
          })
        )
      }

      // zoomreset button (set the image to the "original" size, same size as when selecting the image
      if (options.buttonZoomreset !== false) {
        $(tools).append(
          $(_self.button.zoomreset).on({
            'touchstart click': function (e) {
              e.preventDefault()
              _self.imageReset()
            }
          })
        )
      }

      if (options.buttonFullscreen !== false) {
        $(tools).append(
          $(_self.button.fullscreen).on({
            'touchstart click': function (e) {
              e.preventDefault()
              _self.expandAndResize()
            }
          })
        )
      }

      // zoomout button
      if (options.buttonZoomout !== false) {
        $(tools).append(
          $(_self.button.zoomout).on({
            'touchstart mousedown': function (e) {
              e.preventDefault()
              var imageZoom = 0.5
              var zooming = 0
              _self.interval = window.setInterval(function () {
                zooming++
                if (zooming % 150 === 0) {
                  imageZoom = imageZoom + 0.5
                }
                _self.imageZoom((imageZoom * -1))
              }, 1)
            },
            'touchend mouseup mouseleave': function (e) {
              e.preventDefault()
              window.clearInterval(_self.interval)
              if (_self.options.onAfterZoomImage) { _self.options.onAfterZoomImage.call(_self, _self) }
            }
          })
        )
      }
      // if (options.buttonRotateccw !== false) {
      //   $(tools).append(
      //     $(_self.button.rotateccw).on({
      //       'touchstart touchend click': function (e) {
      //         e.preventDefault()
      //         _self.rotate(-90)
      //       }
      //     })
      //   )
      // }
      // if (options.buttonRotatecw !== false) {
      //   $(tools).append(
      //     $(_self.button.rotatecw).on({
      //       'touchstart touchend click': function (e) {
      //         e.preventDefault()
      //         _self.rotate(90)
      //       }
      //     })
      //   )
      // }
      // cancel button (removes the image and resets it to the original init event
      if (options.buttonCancel !== false) {
        $(tools).append(
          $(_self.button.cancel).on({
            'touchstart touchend click': function (e) {
              e.preventDefault()
              _self.reset()
            }
          })
        )
      }
      // done button (crop the image!)
      if (options.buttonDone !== false) {
        $(tools).append(
          $(_self.button.done).on({
            'touchstart click': function (e) {
              e.preventDefault()
              _self.imageCrop()
            }
          })
        )
      }

      $(element).append($(tools))
    },
    _clearTimers: function () {
      // function to clear all timers, just to be sure!
      var intervalID = window.setInterval(function () {}, 9999)
      for (var i = 1; i < intervalID; i++) { window.clearInterval(i) }
    },
    _formValidation: function () {
      var _self = this
      var element = _self.element
      // var input = _self.input

      $(element).closest('form').submit(function (e) {
        // e.stopPropagation();
        $(this).find('input[type=file]').each(function (i, el) {
          if ($(el).prop('required')) {
            if ($(el).data('valid') !== true) {
              e.preventDefault()
              return false
            }
          }
        })

        return true
      })
    }
  }

  $.fn.ImageUpload = function (options) {
    if ($.data(this, 'ImageUpload')) return
    return $(this).each(function () {
      $.data(this, 'ImageUpload', new $.ImageUpload(options, this))
    })
  }
})(window, jQuery)

function empty (mixedVar) {
  // discuss at: http://phpjs.org/functions/empty/
  // original by: Philippe Baumann
  //    input by: Onno Marsman
  //    input by: LH
  //    input by: Stoyan Kyosev (http://www.svest.org/)
  // bugfixed by: Kevin van Zonneveld (http://kevin.vanzonneveld.net)
  // improved by: Onno Marsman
  // improved by: Francesco
  // improved by: Marc Jansen
  // improved by: Rafal Kukawski

  var undef, key, i, len
  var emptyValues = [undef, null, false, 0, '', '0']

  for ((i = 0), (len = emptyValues.length); i < len; i++) {
    if (mixedVar === emptyValues[i]) {
      return true
    }
  }

  if (typeof mixedVar === 'object') {
    for (key in mixedVar) {
      return false
    }
    return true
  }
  return false
}
