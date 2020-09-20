Vue.component("node-ip-addresses-box", {
	props: ["vIpAddresses"],
	data: function () {
		return {
			ipAddresses: (this.vIpAddresses == null) ? [] : this.vIpAddresses
		}
	},
	methods: {
		// 添加IP地址
		addIPAddress: function () {
			let that = this;
			teaweb.popup("/nodes/ipAddresses/createPopup", {
				callback: function (resp) {
					that.ipAddresses.push(resp.data.ipAddress);
				}
			})
		},

		// 修改地址
		updateIPAddress: function (index, address) {
			let that = this;
			teaweb.popup("/nodes/ipAddresses/updatePopup?addressId=" + address.id, {
				callback: function (resp) {
					Vue.set(that.ipAddresses, index, resp.data.ipAddress);
				}
			})
		},

		// 删除IP地址
		removeIPAddress: function (index) {
			this.ipAddresses.$remove(index);
		}
	},
	template: `<div>
	<input type="hidden" name="ipAddresses" :value="JSON.stringify(ipAddresses)"/>
	<div v-if="ipAddresses.length > 0">
		<div>
			<div v-for="(address, index) in ipAddresses" class="ui label small">
				{{address.ip}}<span class="small" v-if="address.name.length > 0">（{{address.name}}）</span>
				<a href="" title="修改" @click.prevent="updateIPAddress(index, address)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeIPAddress(index)"><i class="icon remove"></i></a>
			</div>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button small" type="button" @click.prevent="addIPAddress()">+</button>
	</div>
	<p class="comment">添加已经绑定的IP地址，仅做记录用。</p>
</div>`
})