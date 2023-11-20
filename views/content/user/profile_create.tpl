<div class="right-col" role="main">
  <div>
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
      <div class="x_title">
        <h2>Create New Profile</h2>
      </div>
      <div class="x_content">
        <div class="well">
          Fill out the form below to create a new profile for User: <b>{{ .userEmail }}</b>
        </div>
        <form method="POST" action="/users/{{ .userID }}/profile_add" id="demo-form2" class="form-horizontal">
          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12"><span class="required">*</span>Description</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <input id="description" name="description" type="text" class="form-control" placeholder="Description" data-parsley-maxlength="255" data-parsley-required="true" data-parsley-trigger="change" required>
            </div>
          </div>

          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12">Delivery Method</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <div class="checkbox">
                <label class="label-cr">
                  {{ if eq .CanSendEmail false }}
                  <input id="delivery" name="delivery" type="checkbox" value="0" disabled>
                  {{ else }}
                  <input id="delivery" name="delivery" type="checkbox" value="1" checked="checked">
                  {{ end }}
                  <span class="labelbox">
                    <i class="ico"></i>Send VPN profile to user via email
                    {{ if eq .CanSendEmail false }}
                    <i class="ion ion-information-circled" data-toggle="tooltip" data-placement="right" title="" data-original-title="Go to 'Settings' > 'Mail' to enable this feature" style="padding-left:5px"></i>
                    {{ end }}
                  </span>
                </label>
              </div>
            </div>
          </div>

          <div class="ln_solid"></div>

          <div class="form-group pull-right">
            <div class="col-md-12">
              <button type="submit" class="btn btn-primary" data-disable-with="Create Profile">Create Profile</button>
              <a href="/users/{{ .userID }}" class="btn btn-default">Cancel</a>
            </div>
          </div>
        </form>
      </div>
    </div>

    <div class="clearfix"></div>
    <br />
  </div>
</div>
