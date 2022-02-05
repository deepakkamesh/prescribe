
// Control Websocket message handlers
$(document).ready(function() {
  var errorContainer = document.querySelector('#error-popup');
});

// Configuration Control.
$(document).ready(function() {
  document.querySelector('#master_enable')
      .addEventListener('click', function() {
        if (document.getElementById('master_enable').checked) {
          SendControlCmd(CmdType.MASTER_ENABLE);
        } else {
          SendControlCmd(CmdType.MASTER_DISABLE);
        }
      });

  document.querySelector('#audio_enable').addEventListener('click', function() {
    if (document.getElementById('audio_enable').checked) {
      SendControlCmd(CmdType.AUDIO_ENABLE);
    } else {
      SendControlCmd(CmdType.AUDIO_DISABLE);
    }
  });

  document.querySelector('#headlight_enable')
      .addEventListener('click', function() {
        if (document.getElementById('headlight_enable').checked) {
          SendControlCmd(CmdType.HEADLIGHT_ON);
        } else {
          SendControlCmd(CmdType.HEADLIGHT_OFF);
        }
      });

  document.querySelector('#video_enable').addEventListener('click', function() {
    fps = parseInt($('#fps_sel').val());
    resMode = parseInt($('#res-sel').val());
    data = [ fps, resMode ];
    if (document.getElementById('video_enable').checked) {
      SendControlCmd(CmdType.VIDEO_ENABLE, data);
      $("#video_stream")
          .attr("src", "/videostream" +
                           '?' + Math.random());
    } else {
      $("#video_stream").attr("src", "");
      SendControlCmd(CmdType.VIDEO_DISABLE, data);
    }
  });
});

// Servo and Drive Controls.
$(document).ready(function() {
  document.querySelector('#generate_pdf').addEventListener('click', function() {
    name = $('#name').val();
    age_sex = $('#age_sex').val();
    prescription = $('#prescription').val();

    $.post('/api/genpdf',
           {name : name, age_sex : age_sex, prescription : prescription},
           function(data, status) {
             console.log(data.Data);
             if (data.Err != '') {
               console.log(data.Err);
               return
             }
           });
  });

  // TODO: DEL Motor Controls.
  document.querySelector('#motor-forward')
      .addEventListener('click', function() {
        SendControlCmd(CmdType.DRIVE_FWD,
                       parseInt($('#drive_velocity_sel').val()));
      });

  document.querySelector('#motor-back').addEventListener('click', function() {
    SendControlCmd(CmdType.DRIVE_BWD, parseInt($('#drive_velocity_sel').val()));
  });

  document.querySelector('#motor-right').addEventListener('click', function() {
    if (document.getElementById('rotate_dual').checked) {
      SendControlCmd(CmdType.DRIVE_RIGHT,
                     parseInt($('#drive_velocity_sel').val()));
      return;
    }
    SendControlCmd(CmdType.DRIVE_RIGHT_ONLY,
                   parseInt($('#drive_velocity_sel').val()));
  });

  document.querySelector('#motor-left').addEventListener('click', function() {
    if (document.getElementById('rotate_dual').checked) {
      SendControlCmd(CmdType.DRIVE_LEFT,
                     parseInt($('#drive_velocity_sel').val()));
      return;
    }
    SendControlCmd(CmdType.DRIVE_LEFT_ONLY,
                   parseInt($('#drive_velocity_sel').val()));
  });

  // Drive velocity selector.
  document.querySelector('#drive_velocity_sel')
      .addEventListener('click', function() {
        val = $('#drive_velocity_sel').val();
        $("#drive_velocity_sel_disp").empty()
        $("#drive_velocity_sel_disp").append(val);
      });

  // Servo Controls.
  document.querySelector('#servo-down')
      .addEventListener('click',
                        function() { SendControlCmd(CmdType.SERVO_DOWN); });

  document.querySelector('#servo-up')
      .addEventListener('click',
                        function() { SendControlCmd(CmdType.SERVO_UP); });

  document.querySelector('#servo-top')
      .addEventListener('click',
                        function() { SendControlCmd(CmdType.SERVO_ABS, 0); });

  document.querySelector('#servo-center')
      .addEventListener('click',
                        function() { SendControlCmd(CmdType.SERVO_ABS, 90); });

  document.querySelector('#servo-bottom')
      .addEventListener('click',
                        function() { SendControlCmd(CmdType.SERVO_ABS, 180); });

  // Set the step for Servo.
  document.querySelector('#servo_angle_step')
      .addEventListener('click', function() {
        val = $('#servo_angle_step').val();
        $("#servo_angle_step_disp").empty();
        $("#servo_angle_step_disp").append(val);

        SendControlCmd(CmdType.SERVO_STEP, parseInt(val));
      });
});
