<div class="top-nav">
  <div class="nav-menu">
    <nav>
      <div class="nav toggle">
        <a id="menu_toggle"><i class="bars"></i></a>
      </div>

      <ul class="nav navbar-nav navbar-right">
        <li>
          <a href="/sign_out"> <span>Log out</span> &ensp;<i class="ion ion-android-exit"></i></a>
        </li>
        <li>
          <a href="{{.UserLink}}" class=""> <span> {{.Username}} </span> <small> {{.Role}} </small></a>
        </li>
      </ul>
    </nav>
  </div>
</div>