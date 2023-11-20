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
          <h1>Sign In</h1>
        </div>

        <div class="x_content">
          <form method="POST" id="demo-form2" data-parsley-validate class="form-horizontal">
            <div class="form-group">
              <label class="col-md-12 col-sm-12 col-xs-12" for="email">Email</label>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="email" id="email" name="email" required class="form-control" placeholder="Email" value="{{ .Form.Email }}" />
              </div>
            </div>
            <div class="form-group">
              <label class="col-md-12 col-sm-12 col-xs-12" for="password">Password</label>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="password" id="password" name="password" required class="form-control" placeholder="Password" value="" />
              </div>
            </div>
            <button type="submit" class="btn btn-block btn-primary" style="margin: 10px 0 10px 0" data-disable-with="Sign In">Sign In</button>
          </form>
          <a href="/password_recovery" class="btn btn-block btn-default" id="forgotPassword">Forgot Password</a>
        </div>
      </div>
    </div>
  </div>
</div>
