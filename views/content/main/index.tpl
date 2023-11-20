<div class="right-col" role="main">
  <div>
    <div class="clearfix"></div>

    <div class="row tile_count">
      <div class="col-md-3 col-sm-6 col-xs-6 tile_stats_count">
        <span class="count_top"><i class="ion ion-android-time"></i> Uptime </span>
        <div class="count">{{ .UpTime }}</div>
        <!-- <span class="count_bottom"><i class="green">4% </i> From last Week</span> -->
      </div>
      <div class="col-md-3 col-sm-6 col-xs-6 tile_stats_count">
        <span class="count_top"><i class="ion ion-android-bulb"></i> Availability </span> {{ if eq .Availability 0 }}
        <div class="count text-center"> <i class="glyphicon glyphicon-ok-sign" style="color: #0E0; border-radius: 16px; background-color: #FFF;"></i> </div>
        {{ else }}
        <div class="count text-center"> <i class="glyphicon glyphicon-remove-sign" style="color: #E00; border-radius: 16px; background-color: #FFF;"></i> </div>
        {{ end }}
      </div>
      <div class="col-md-3 col-sm-6 col-xs-6 tile_stats_count">
        <span class="count_top"><i class="ion ion-log-in"></i> Total Incoming Traffic </span>
        <div class="count"> {{ .TotalTrafficIncoming }} </div>
      </div>
      <div class="col-md-3 col-sm-6 col-xs-6 tile_stats_count">
        <span class="count_top"><i class="ion ion-log-out"></i> Total Outgoing Traffic </span>
        <div class="count"> {{ .TotalTrafficOutgoing }} </div>
      </div>
    </div>
    <div class="clearfix"></div>
    <!-- /top tiles -->

    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
        <div class="x_panel">
          <div class="x_title">
            <h2>Currently Active Sessions</h2>
            <div class="clearfix"></div>
          </div>
          <div class="x_content">
            <div id="map" style="height: 400px;width: 100%;"></div>
          </div>
        </div>
      </div>
    </div>

    <div class="clearfix"></div>

    <div class="row">
      <div class="col-md-12 col-xs-12">
        <div class="x_panel">
          <div class="x_title">
            <h2>Session List</h2>
            <div class="clearfix"></div>
          </div>
          <div class="x_content">
            <div class="table-responsive">
              <table class="table table-striped">
                <thead>
                  <tr>
                    <th>User Name</th>
                    <th>Description</th>
                    <th>Session Name</th>
                    <th>Source IP</th>
                    <th>Session Start</th>
                    <th>Download / Upload</th>
                  </tr>
                </thead>
                <tbody>
                  {{ if .DataList }} {{ range .DataList }}
                  <tr>
                    <th scope="row"><a href="/users/{{.UserID}}"> {{.UserEmail}} </a></th>
                    <td>{{ .UserAlias }}</td>
                    <td>{{ .SessionName }}</td>
                    <td>{{ .SourceIP }}</td>
                    <td>{{ .SessionStart }}</td>
                    <td>
                      <i class="ion ion-ios-cloud-download-outline"></i> {{ .IncomingByte }}
                      <span>/</span>
                      <i class="ion ion-ios-cloud-upload-outline"></i> {{ .OutgoingByte }}
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

  </div>
</div>