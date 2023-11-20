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
        <h1>Logs</h1>
      </div>
    </div>

    <div class="clearfix"></div>
    <hr>

    <div class="clearfix"></div>
    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
        <div class="form-group pull-right">
          <a href="/logs/download" class="btn btn-primary " id="exportLog"><i class="ion ion-archive"></i> Download as CSV</a>
        </div>
        <div class="x_panel">
          <div class="x_content">
            <div class="table-responsive">
              <table class="table table-striped">
                <thead>
                <tr>
                  <th>Time</th>
                  <th>Type</th>
                  <th>Message</th>
                </tr>
                </thead>
                <tbody>
                {{ if .LogList }} {{ range .LogList }}
                <tr id="tr_{{ .Id }}" class="clickable subspace-list-hover" data-toggle="collapse" data-target="#tr_target_{{ .Id }}">
                  <td>{{ localTimeFmt .LogTime }}</td>
                  <td>{{ .Type }}</td>
                  <td>{{ substr .RawLog 0 80 }} {{ $length := len .RawLog }} {{ if gt $length 80 }} ... {{ end }}
                  </td>
                </tr>
                <tr>
                  <td colspan="8">
                    <div class="accordion-body collapse" id="tr_target_{{ .Id }}">
                      <div class="beautify-json">{{ .RawLog }}</div>
                    </div>
                  </td>
                </tr>
                {{ end }} {{ end }}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!---->

  </div>
</div>
