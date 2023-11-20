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

    <!-- Main Content -->
    <div class="page-title">
      <div class="title_left">
        <h1>Backup / Restore</h1>
      </div>
    </div>
    <!-- Host info -->
    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
        <div class="x_panel">
          <div class="x_title">
            <h2>Backup <small>latest backup time: {{ localTimeFmt .SystemInfo.LatestBackupDate }}</small></h2>
            <div class="clearfix"></div>
          </div>
          <div class="x_content">
            <form class="form-horizontal">
              <div class="form-group">
                <div class="col-md-12">
                  <p>
                    Download Backup File.
                  </p>
                  <div id="finish_backup" class="alert alert-success alert-dismissible fade in" role="alert" style="display:none;">
                    <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">×</span>
                    </button>
                    <span id="backup_completed_msg"></span>
                  </div>
                  <div id="failed_backup" class="alert alert-danger alert-dismissible fade in" role="alert" style="display:none;">
                    <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">×</span>
                    </button>
                    <span id="backup_failed_msg"></span>
                  </div>
                  <button id="start_backup" class="btn btn-primary" type="button"><i class="ion ion-archive" style="font-size:16px; padding-right:8px"></i> Backup </button>
                  <button id="backing_up" class="btn btn-danger" type="button" style="display:none;" disabled><i class="fa fa-spinner fa-pulse fa-lg fa-fw"></i> Backing Up... </button>
                  <a id="download_backup" class="btn btn-success" href="" style="display:none;" download><i class="ion ion-archive" style="font-size:16px; padding-right:8px"></i> Download </a>
                  <button id="backup_again" class="btn btn-primary" type="button" style="display:none;"><i class="ion ion-archive" style="font-size:16px; padding-right:8px"></i> Backup Again </button>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    <!-- /Host info -->

    <div class="ln_solid"></div>
    <div class="clearfix"></div>

    <!-- VPN info -->
    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
        <div class="x_panel">
          <div class="x_title">
            <h2>Restore</h2>
            <div class="clearfix"></div>
          </div>
          <div class="x_content">
            <form id="restore" method="POST" action="" class="form-horizontal" enctype="multipart/form-data">
              <div class="form-group">
                <div class="col-md-12">
                  <p>
                    Upload configuration file to restore service.
                  </p>
                  <input type="file" name="upload_config" id="upload_config" class="" title="Browse file"/>
                </div>
              </div>
              <!--<br>-->
              <div class="form-group">
                <div class="col-md-12">
                  <div id="finish_restore" class="alert alert-success alert-dismissible subspace-alert-style fade in" role="alert" style="display:none;">
                    <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">×</span>
                    </button>
                    <span id="restore_completed_msg"></span>
                    <!--Restore process using config: <span id="config_filename"></span> is accomplished, proceed to use VPN service.-->
                  </div>
                  <div id="failed_restore" class="alert alert-danger alert-dismissible subspace-alert-style fade in" role="alert" style="display:none;">
                    <button type="button" class="close" data-dismiss="alert" aria-label="Close"><!-- <span aria-hidden="true">×</span> -->
                    </button>
                    <span id="restore_failed_msg">Restore unsuccessful. Please choose another file to try again or contact support.</span>
                  </div>
                  <button id="restore_again" type="button" form="restore" class="btn btn-primary " style="display:none;"><i class="ion ion-archive" style="font-size:16px; padding-right:8px"></i> Restore again </button>
                </div>
              </div>
              <div class="form-group">
                <div class="col-md-12">
                  <button id="start_restore" type="button" form="restore" class="btn btn-primary "><i class="ion ion ion-archive" style="font-size:16px; padding-right:8px"></i> Restore</button>
                  <button id="restoring" class="btn btn-danger" type="button" style="display:none;" disabled><i class="fa fa-spinner fa-pulse fa-lg fa-fw"></i>  Restoring... </button>
                  <p id="restore_message" style="color:red;margin-top:10px;"></p>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    <!-- /VPN info -->
    <!-- /Main Content -->
  </div>
</div>

<div id="myNav" class="overlay">
  <!-- Button to close the overlay navigation -->
  <!--<a href="javascript:void(0)" class="closebtn" onclick="closeNav()">&times;</a>-->
  <!-- Overlay content -->
  <div class="overlay-content">
    <a href="#"><i class="fa fa-spinner fa-pulse fa-lg fa-fw"></i>  Restoring...</a>
    <br>
    <p>Please wait while the your settings are restored. The whole process may take a few minutes.</p>
    <p>You will be redirected to the log in page after the restoration process is complete. Thank you for your patience.</p>
  </div>
</div>
