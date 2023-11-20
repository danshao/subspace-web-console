<div>
  <div class="row">
    <div class="col-md-8 col-sm-8 col-xs-12 col-md-offset-2 col-sm-offset-2">
      <div class="x_panel">
        <div class="x_title subspace-h1-title-visible">
          <h1>Welcome to Subspace!</h1>
        </div>

        <div class="x_content">
          <form method="POST" action="" id="onboardStartForm" data-parsley-validate class="form-horizontal">
            <div class="form-group">
              <label for="instanceID" class="col-md-12 col-sm-12 col-xs-12">
                EC2 Instance ID
                <span class="required">*</span>
              </label>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="text" id="instanceID" name="instanceID" required="required" class="form-control" placeholder="Enter Your EC2 Instance ID">
              </div>
            </div>
            <div>
              {{if .flash.error}}
              <font color="red">{{.flash.error}}</font>
              {{end}}
            </div>
            <div class="form-help">
              <p class="help-block">Your Subspace EC2 instance ID can be found in the EC2 Management Area of the AWS Console.</p>
            </div>
            <div class="setup-action-Btn">
              <button type="submit" class="btn btn-success btn-block" data-disable-with="Continue">Continue</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</div>
