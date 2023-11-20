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
        <h1>User Management</h1>
      </div>
    </div>
    <hr>

    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
          <a class="btn btn-primary pull-right" style="margin-bottom: 10px" href="/users/add" role="button"> <i class="ion ion-android-person-add"></i> Add User </a>
      </div>
    </div>

    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
        <div class="x_panel">
          <div class="x_title">
            <h2>Subspace Users</h2>
          </div>
          <div class="x_content">
            <div class="table-responsive">
              <table class="table table-striped">
                <thead>
                  <tr>
                    <th>Email</th>
                    <th>Alias</th>
                    <th>Role</th>
                    <th>Status</th>
                  </tr>
                </thead>
                {{ if .UserList}}
                <tbody>
                  {{ range .UserList }}
                  <tr>
                    <th scope="row"><a href="/users/{{ .Id }}">{{ .Email }}</a></th>
                    <td>{{ .Alias }}</td>
                    <td>{{ if eq .Role "admin" }}Administrator
                        {{ else }}User
                        {{ end }}
                    </td>
                    <td>{{ if eq .Enabled true }}Enabled
                        {{ else }}Disabled
                        {{ end }}
                    </td>
                  </tr>
                  {{ end }}
                </tbody>
                {{ end }}
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
