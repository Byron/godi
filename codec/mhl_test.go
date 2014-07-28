package codec

import (
	"encoding/xml"
	"testing"

	"github.com/Byron/godi/api"
)

func TestMHLDecode(t *testing.T) {
	hl := mhlHashList{}
	if err := xml.Unmarshal([]byte(mhlFixture), &hl); err != nil {
		t.Fatal(err)
	}

	if len(hl.HashInfo) != 16 {
		t.Errorf("Expected 16 entries, got %d", len(hl.HashInfo))
	}

	fi := api.FileInfo{}
	for _, h := range hl.HashInfo {
		if err := h.ToFileInfo(&fi); err != nil {
			t.Error(err)
		}
	}
}

const mhlFixture = `<?xml version="1.0" encoding="UTF-8"?>
<hashlist version="1.0">

  <creatorinfo>
    <name>Firname Lastname</name>
    <username>login</username>
    <hostname>hostname.localdomain</hostname>
    <tool>mhl ver. 0.1.28</tool>
    <startdate>2014-07-28T14:17:57Z</startdate>
    <finishdate>2014-07-28T14:19:09Z</finishdate>
    <log><![CDATA[===================
mhl ver. 0.1.28 started.
Verbose mode is ON.
Working directory: "..."
MHL file directory(-es):
   .../uhd-demo
-------------------
-------------------
Start date in UTC: 2014-07-28 14:17:57.
-------------------
Calculating hash sums
Processing '.../uhd-demo/UHD_demo_a.mp4'
Done '.../uhd-demo/UHD_demo_a.mp4'
Processing '.../uhd-demo/UHD_demo_a_Astra.mp4'
Done '.../uhd-demo/UHD_demo_a_Astra.mp4'
Processing '.../uhd-demo/UHD_demo_a_DEcityscape.mp4'
Done '.../uhd-demo/UHD_demo_a_DEcityscape.mp4'
Processing '.../uhd-demo/UHD_demo_a_Dustin_Germany.mp4'
Done '.../uhd-demo/UHD_demo_a_Dustin_Germany.mp4'
Processing '.../uhd-demo/UHD_demo_b.mp4sub'
Done '.../uhd-demo/UHD_demo_b.mp4sub'
Processing '.../uhd-demo/UHD_demo_b_Astra.mp4sub'
Done '.../uhd-demo/UHD_demo_b_Astra.mp4sub'
Processing '.../uhd-demo/UHD_demo_b_DEcityscape.mp4sub'
Done '.../uhd-demo/UHD_demo_b_DEcityscape.mp4sub'
Processing '.../uhd-demo/UHD_demo_b_Dustin_Germany.mp4sub'
Done '.../uhd-demo/UHD_demo_b_Dustin_Germany.mp4sub'
Processing '.../uhd-demo/UHD_demo_c.mp4sub'
Done '.../uhd-demo/UHD_demo_c.mp4sub'
Processing '.../uhd-demo/UHD_demo_c_Astra.mp4sub'
Done '.../uhd-demo/UHD_demo_c_Astra.mp4sub'
Processing '.../uhd-demo/UHD_demo_c_DEcityscape.mp4sub'
Done '.../uhd-demo/UHD_demo_c_DEcityscape.mp4sub'
Processing '.../uhd-demo/UHD_demo_c_Dustin_Germany.mp4sub'
Done '.../uhd-demo/UHD_demo_c_Dustin_Germany.mp4sub'
Processing '.../uhd-demo/UHD_demo_d.mp4sub'
Done '.../uhd-demo/UHD_demo_d.mp4sub'
Processing '.../uhd-demo/UHD_demo_d_Astra.mp4sub'
Done '.../uhd-demo/UHD_demo_d_Astra.mp4sub'
Processing '.../uhd-demo/UHD_demo_d_DEcityscape.mp4sub'
Done '.../uhd-demo/UHD_demo_d_DEcityscape.mp4sub'
Processing '.../uhd-demo/UHD_demo_d_Dustin_Germany.mp4sub'
Done '.../uhd-demo/UHD_demo_d_Dustin_Germany.mp4sub'
-------------------
End of input.
Finish date in UTC: 2014-07-28 14:19:09
MHL file path(s):
   .../uhd-demo/uhd-demo_2014-07-28_141757.mhl
===================
]]>
    </log>
  </creatorinfo>

  <hash>
    <file>UHD_demo_a.mp4</file>
    <size>1119644340</size>
    <lastmodificationdate>2013-09-06T07:47:22Z</lastmodificationdate>
    <sha1>2aed85b910adfd0a74ef93d50c802a311c7ef710</sha1>
    <md5>cd25f7f81ee1780508f6c98a3eaea1f3</md5>
    <hashdate>2014-07-28T14:18:02Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_a_Astra.mp4</file>
    <size>1110631034</size>
    <lastmodificationdate>2013-09-06T09:07:38Z</lastmodificationdate>
    <sha1>f2d2eb213c1162956c5f72b6a7751c092df55a02</sha1>
    <md5>f313dcf60adaa0d426a3bf08f2c632f9</md5>
    <hashdate>2014-07-28T14:18:08Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_a_DEcityscape.mp4</file>
    <size>369962740</size>
    <lastmodificationdate>2013-08-10T04:22:56Z</lastmodificationdate>
    <sha1>2f6fe489fb1a221280cfc58a0e9572ee2ddf511d</sha1>
    <md5>59f89a9a4d1df85bf82c5eefde7515bd</md5>
    <hashdate>2014-07-28T14:18:10Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_a_Dustin_Germany.mp4</file>
    <size>478842307</size>
    <lastmodificationdate>2013-08-10T04:16:46Z</lastmodificationdate>
    <sha1>fb4adfb1e670533eadf6c593139698e8eb09ae07</sha1>
    <md5>6da66a14957d724f2c08ba14f60b387e</md5>
    <hashdate>2014-07-28T14:18:13Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_b.mp4sub</file>
    <size>1119634188</size>
    <lastmodificationdate>2013-09-05T20:23:56Z</lastmodificationdate>
    <sha1>a736f4e1c041f1975c71e21fe1b12f01788f98ba</sha1>
    <md5>3c575a5e911e99525652a120472b43e9</md5>
    <hashdate>2014-07-28T14:18:20Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_b_Astra.mp4sub</file>
    <size>1124461013</size>
    <lastmodificationdate>2013-09-05T20:25:42Z</lastmodificationdate>
    <sha1>70a9343d21c2c16d1e9a866472ccc5ac883035ac</sha1>
    <md5>b4afc752cd298bd9a03c07390d103998</md5>
    <hashdate>2014-07-28T14:18:27Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_b_DEcityscape.mp4sub</file>
    <size>380419644</size>
    <lastmodificationdate>2013-08-10T04:23:02Z</lastmodificationdate>
    <sha1>0b7560e36eda8bfe14dbb64d60c331b1abc48d3f</sha1>
    <md5>0f719c607b1bada75216c93076144e9b</md5>
    <hashdate>2014-07-28T14:18:29Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_b_Dustin_Germany.mp4sub</file>
    <size>478842439</size>
    <lastmodificationdate>2013-08-10T04:19:24Z</lastmodificationdate>
    <sha1>76739c1baeccd341437200c1b76c1e8764e80b6b</sha1>
    <md5>7560440a2f8170eabaab7aa1e7fdee64</md5>
    <hashdate>2014-07-28T14:18:32Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_c.mp4sub</file>
    <size>1119640204</size>
    <lastmodificationdate>2013-09-05T20:26:18Z</lastmodificationdate>
    <sha1>5b54e928377e4371585330e47af10b827390d832</sha1>
    <md5>c7b10e235e0f9284c28d10a233af631c</md5>
    <hashdate>2014-07-28T14:18:38Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_c_Astra.mp4sub</file>
    <size>1098055344</size>
    <lastmodificationdate>2013-09-05T20:21:50Z</lastmodificationdate>
    <sha1>416b1685ec155066c2295dfbdcd1e5c76d2b6498</sha1>
    <md5>7fe22c23c7f90db21394e394e4ca2cbc</md5>
    <hashdate>2014-07-28T14:18:45Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_c_DEcityscape.mp4sub</file>
    <size>368527981</size>
    <lastmodificationdate>2013-08-10T04:23:00Z</lastmodificationdate>
    <sha1>db4f289dc0076c477561e3503f74166eb1255aeb</sha1>
    <md5>a83d394a779278c7525fd17a27062ed6</md5>
    <hashdate>2014-07-28T14:18:47Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_c_Dustin_Germany.mp4sub</file>
    <size>478842247</size>
    <lastmodificationdate>2013-08-10T04:20:06Z</lastmodificationdate>
    <sha1>0a9e172f700196d6757e8ef93597f7b6a42a90eb</sha1>
    <md5>15feec4ad91cd516db546e80c49a01d0</md5>
    <hashdate>2014-07-28T14:18:50Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_d.mp4sub</file>
    <size>1119634188</size>
    <lastmodificationdate>2013-09-05T20:26:04Z</lastmodificationdate>
    <sha1>83177ac75966d074e5fdc9d5f6ec28b10af71b13</sha1>
    <md5>49907cdd7d00b47ab78bcadc92e7be4d</md5>
    <hashdate>2014-07-28T14:18:57Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_d_Astra.mp4sub</file>
    <size>1117435755</size>
    <lastmodificationdate>2013-09-05T20:09:36Z</lastmodificationdate>
    <sha1>87f5d314534562db0f67e63ad596f2982c2a7ed8</sha1>
    <md5>b752e76b7f8f5e548ec15eeb03a56e6c</md5>
    <hashdate>2014-07-28T14:19:04Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_d_DEcityscape.mp4sub</file>
    <size>381182324</size>
    <lastmodificationdate>2013-08-10T04:23:06Z</lastmodificationdate>
    <sha1>561e4bb51dac16a356c081a01add8068dd764707</sha1>
    <md5>8aebf8170ec08c52e8868f76ec6e8127</md5>
    <hashdate>2014-07-28T14:19:06Z</hashdate>
  </hash>

  <hash>
    <file>UHD_demo_d_Dustin_Germany.mp4sub</file>
    <size>478842307</size>
    <lastmodificationdate>2013-08-10T04:20:00Z</lastmodificationdate>
    <sha1>30f0dc447a1073a087862fa4bceed3584c0e4834</sha1>
    <md5>a2a3c0a996c5320c1482ad80c03fadd6</md5>
    <hashdate>2014-07-28T14:19:09Z</hashdate>
  </hash>

</hashlist>
`
