<!DOCTYPE html>
<html lang="en">
  <head>
    <!-- header -->
    {{.Header}}
    <!-- /header -->
  </head>

  <body class="nav-md">
    <div class="container body">
      <div class="main-container">

        <!-- sidebar -->
        {{.Sidebar}}
        <!-- /sidebar -->

        <!-- top navigation -->
        {{.Nav}}
        <!-- /top navigation -->

        <!-- page content -->
        {{.LayoutContent}}
        <!-- /page content -->

        <!-- footer content -->
        {{.Footer}}
        <!-- /footer content -->
      </div>
    </div>

    <!-- scripts -->
    {{.Scripts}}
    {{.Scripts_Google_Map}}
    <!-- /scripts -->
  </body>
</html>
