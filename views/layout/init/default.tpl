<!DOCTYPE html>
<html lang="en">

<head>
  <!-- header -->
  {{.Header}}
  <!-- /header -->
</head>

<body class="login">
  <div class="container-fluid body subspace-setup-margin-footer">
    <div class="main-container">
      <div class="navbar nav-title" style="background: none; margin: 20px auto 60px auto; width: 300px; height: 70px;">
        <div class="site_title" style="margin: 0 auto;">
          <img src="/frontend/build/images/logo_badge.svg" class="logo-sm" alt="Subspace">   
          <img src="/frontend/build/images/logo.svg" class="logo" alt="Subspace">  
        </div>
      </div>
      <!-- content -->
      {{.LayoutContent}}
      <!-- /content -->

    </div>
  </div>

  <!-- footer -->
  {{.Footer}}
  <!-- /footer -->

  <!-- script -->
  {{.Scripts}} {{.Scripts_Custom}}
  <!-- /script -->
</body>

</html>