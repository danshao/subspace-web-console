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

    <div class="page-title">
      <div class="title_left">
        <h1>About</h1>
      </div>
    </div>

    <div class="clearfix"></div>
    <hr>

    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
        <div class="x_panel">
          <div class="x_content">
            <h4 style="font-size:20px">Subspace</h4>
            A modern VPN service

            <div style="padding-top:60px">
              Need help? Visit the <a href="http://support.ecowork.com">Support Center</a>.
            </div>

            <div style="padding-top:25px">
              Version {{.VersionNumber}} Build {{.BuildNumber}}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
