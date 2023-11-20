<div>
  <div class="row">
    <div class="col-md-12 col-sm-12 col-xs-12 single">
      <div class="clearfix"></div>

      <!-- flash message -->
      {{ if .flash.notice }}
      <div class="alert alert-success alert-dismissible fade in" role="alert">
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
        {{ .flash.notice }}
      </div>
      {{ end }}

      {{ if .flash.error }}
      <div class="alert alert-danger alert-dismissible fade in" role="alert">
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
        {{ .flash.error }}
      </div>
      {{ end }}

      {{ if .flash.warning }}
      <div class="alert alert-warning alert-dismissible fade in" role="alert">
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
        {{ .flash.warning }}
      </div>
      {{ end }}
      <!-- /flash message -->

      <div class="clearfix"></div>

      <div class="x_panel">
        <div class="x_title subspace-h1-title-visible">
          <h1>Reset Password</h1>
        </div>
        <div class="well" style="margin-bottom: 0px">
          Please fill out the form below to reset your password.
        </div>

        <div class="x_content">
          <br />
          <form method="POST" id="demo-form2" data-parsley-validate class="form-horizontal">
            <div class="form-group">
              <label class="col-md-12 col-sm-12 col-xs-12" for="email"><span class="required">*</span>Administrator Email</label>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="email" id="email" name="email" data-parsley-required="true" data-parsley-trigger="change" class="form-control" placeholder="Email" value="{{ .Form.Email }}" required />
              </div>
            </div>
            <div class="form-group">
              <label class="col-md-12 col-sm-12 col-xs-12" for="instance_id"><span class="required">*</span>Subspace Instance ID</label>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="text" id="instance_id" name="instance_id" data-parsley-required="true" class="form-control" placeholder="Instance ID" value="{{ .Form.InstanceID }}" required />
              </div>
            </div>
            <div class="form-group">
              <label class="col-md-12 col-sm-12 col-xs-12" for="instance_id"><span class="required">*</span>New Password</label>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="password" id="password" name="password" data-parsley-required="true" class="form-control" placeholder="Password" data-parsley-trigger="change" data-parsley-errors-container="#password_error" data-parsley-pattern="(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,}" data-parsley-pattern-message="Your password must contain a minimum of 8 characters and at least 1 upper case letter, 1 lower case letter and 1 number." required />
              </div>
            </div>
            <button type="submit" class="btn btn-block btn-primary" style="margin: 10px 0 10px 0" data-disable-with="Submit">Submit</button>
          </form>
          <a href="/" class="btn btn-block btn-default">Cancel</a>
        </div>
      </div>
    </div>
  </div>
</div>
