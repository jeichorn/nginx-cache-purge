<?php
// this is just a quick test without a testing framework

require __DIR__.'/../vendor/autoload.php';

use NginxCP\Cache;

$keys = 
 array (
     'example.com' => array(
   'httpexample.com/wp-content/plugins/bbpress/templates/default/js/editor.js??2.5.3-5249' => ['/mnt/cache/assets/3/c/66/388893c966a193b41b1b89d7611666c3'],
   'httpexample.com/wp-content/plugins/buddypress/bp-templates/bp-legacy/js/buddypress.js??2.0.1' => ['/mnt/cache/assets/6/9/1a/f83d5a8a5646ca63088c2970f49c1a96'],
   'httpexample.com/wp-includes/js/jquery/jquery-migrate.min.js??1.2.1' => ['/mnt/cache/assets/6/1/b4/0e8dda24513f3da040a4ab8cc90ab416'],
   'httpexample.com/wp-includes/css/buttons.min.css??3.9.1' => ['/mnt/cache/assets/6/d/2c/06b28aec2c42ff301e60b69cfd362cd6'],
   'httpexample.com/wp-includes/css/dashicons.min.css??3.9.1' => ['/mnt/cache/assets/6/4/f5/9016a9994c265ba7111c691abc26f546'],
   'httpexample.com/wp-content/themes/twentytwelve/style.css??3.9.1' => ['/mnt/cache/assets/6/2/25/ce3b5adcd2c0fc1dd557decd26cb2526'],
   'httpexample.com/wp-content/plugins/bbpress/templates/default/css/bbpress.css??2.5.3-5249' => ['/mnt/cache/assets/9/1/e7/7dacbc87ed05e72a680755b2c71fe719'],
   'httpexample.com/wp-includes/js/admin-bar.min.js??3.9.1' => ['/mnt/cache/assets/9/8/02/e8bbaf1a3b71bbfbe169eca18a310289'],
   'httpexample.com/wp-content/themes/twentytwelve/js/navigation.js??20140318' => ['/mnt/cache/assets/b/f/58/8d17d00b1a5614d1dbaed8c28b4a58fb'],
   'httpexample.com/wp-includes/css/editor.min.css??3.9.1' => ['/mnt/cache/assets/8/9/18/06fc913011e7e0e3ffaf3455ca411898'],
   'httpexample.com/wp-includes/js/jquery/jquery.js??1.11.0' => ['/mnt/cache/assets/8/d/bf/fe41f508224498e3419e7bb17160bfd8'],
   'httpexample.com/wp-content/plugins/bbpress/includes/admin/css/admin.css??2.5.3-5249' => ['/mnt/cache/assets/d/0/91/ba1ec695b8b67a4cbda2fa3d5458910d'],
   'httpexample.com/wp-content/plugins/akismet/_inc/form.js??3.0.0' => ['/mnt/cache/assets/c/c/25/50f5ba8d9cddadd4fafee73a2dd525cc'],
   'httpexample.com/wp-content/plugins/buddypress/bp-core/js/confirm.min.js??2.0.1' => ['/mnt/cache/assets/5/6/58/2d099989f901742b1175cde5aa965865'],
   'httpexample.com/wp-content/plugins/buddypress/bp-core/admin/css/common.min.css??2.0.1' => ['/mnt/cache/assets/2/7/dd/1c2180641b6f7527a4cffa7f5329dd72'],
   'httpexample.com/wp-content/plugins/buddypress/bp-templates/bp-legacy/css/buddypress.css??2.0.1' => ['/mnt/cache/assets/2/7/2a/cdeaef483e307e37e9c239f934a32a72'],
   'httpexample.com/wp-includes/js/comment-reply.min.js??3.9.1' => ['/mnt/cache/assets/a/8/6c/8d95ffbfc12d9df9f9efe623eda16c8a'],
   'httpexample.com/wp-content/plugins/buddypress/bp-core/css/admin-bar.min.css??2.0.1' => ['/mnt/cache/assets/e/9/b7/08a51bebf4d828de4b2805b912ceb79e'],
   'httpexample.com/wp-includes/js/thickbox/loadingAnimation.gif?' => ['/mnt/cache/assets/e/7/ed/277e6c8c2de45a00cfed4c8eae74ed7e'],
   'standard--httpsexample.com/wp-includes/css/admin-bar.min.css??3.9.1' => ['/mnt/cache/assets/f/6/60/5eb4ffe9ef449788d53f6a4b7f8a606f'],
   'standard--httpsexample.com/' => ['/mnt/cache/assets/f/6/60/5eb4ffe9ef449788d53f6a4b7f8a606f'],
   'mobile--httpexample.com/wp-content/plugins/akismet/_inc/akismet.css??3.0.0' => ['/mnt/cache/assets/f/6/41/e0dd3e651afc4c72826e95a55cdf416f'],
   'standard--httpexample.com/wp-admin/css/login.min.css??3.9.1' => ['/mnt/cache/assets/f/2/74/9d8754dbf9c9dd904c8a354cf931742f'],
   'httpexample.com/wp-admin/images/wordpress-logo.svg??20131107' => ['/mnt/cache/assets/0/1/8b/1c2e354f5c1f7c5ec893131d35108b10'],
   'httpsexample.com/wp-includes/js/thickbox/thickbox.css??3.9.1' => ['/mnt/cache/assets/0/d/d7/7fb1916bb90ece0f44053d7860c6d7d0'],
   'httpsexample.com/favicon.ico?' => ['/mnt/cache/assets/0/c/2c/1e62c712bcdb67161f88e885fcb72cc0'],
   'httpsexample.com/wp-content/plugins/akismet/_inc/akismet.js??3.0.0' => ['/mnt/cache/assets/0/f/6b/c6213ab04583d93f8a4fdf792c576bf0'],
   'httpsexample.com/foo' => ['/mnt/cache/assets/0/f/6b/c6213ab04583d93f8a4fdf792c576bf1'],
   'httpsexample.com/foo/' => ['/mnt/cache/assets/0/f/6b/c6213ab04583d93f8a4fdf792c576bf2'],
   'httpsexample.com/foo/?' => ['/mnt/cache/assets/0/f/6b/c6213ab04583d93f8a4fdf792c576bf2'],
   'httpsexample.com/foo/?blah=1' => ['/mnt/cache/assets/0/f/6b/c6213ab04583d93f8a4fdf792c576bf2'],
   'example.com/wp-content/plugins/akismet/_inc/akismet.js??3.0.0' => ['/mnt/cache/assets/0/f/6b/c6213ab04583d93f8a4fdf792c576bf0'],
   ),
   'bar.com' => array(
       'httpsbar.com/foo/' => ['/mnt/cache/assets/0/f/6b/c6213ab04583d93f8a4fdf792c576bf2'],
   )
 );

$cases = array(
    'example.com::/(.*)' => count($keys['example.com']),
    'example.com::/' => 1,
    'example.com::/wp-content/(.*)' => 14,
    'example.com::/foo/' => 3,
);

$STATUS = 'OK';
foreach($cases as $case => $num_to_purge)
{
    $purge_rule = $case;
    Cache::$CACHE_PATH = __DIR__;
    $cache = new Cache();
    $cache->keys = $keys;

    $count = $cache->purge($purge_rule);

    if ($count != $num_to_purge)
    {
        echo "$case TEST FAILED\n";
        $STATUS = 'FAILED';
    }
    else
    {
        echo "$case TEST OK\n";
    }
}

echo "Test Status: $STATUS\n";

if ($STATUS == "OK") {
    exit(0);
} else {
    exit(1);
}