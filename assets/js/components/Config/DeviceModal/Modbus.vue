<template>
	<FormRow
		v-if="showConnectionOptions"
		id="modbusTcpIp"
		:label="$t('config.modbus.connection')"
		:help="
			connection === 'tcpip'
				? $t('config.modbus.connectionHintTcpip')
				: $t('config.modbus.connectionHintSerial')
		"
	>
		<div class="btn-group" role="group">
			<input
				id="modbusTcpIp"
				v-model="connection"
				type="radio"
				class="btn-check"
				name="modbusConnection"
				value="tcpip"
				tabindex="0"
				autocomplete="off"
			/>
			<label class="btn btn-outline-primary" for="modbusTcpIp">
				{{ $t("config.modbus.connectionValueTcpip") }}
			</label>
			<input
				id="modbusSerial"
				v-model="connection"
				type="radio"
				class="btn-check"
				name="modbusConnection"
				value="serial"
				tabindex="0"
				autocomplete="off"
			/>
			<label class="btn btn-outline-primary" for="modbusSerial">
				{{ $t("config.modbus.connectionValueSerial") }}
			</label>
		</div>
	</FormRow>
	<FormRow id="modbusId" :label="$t('config.modbus.id')">
		<PropertyField
			id="modbusId"
			property="id"
			type="Int"
			class="me-2"
			required
			:model-value="id || defaultId || 1"
			@change="$emit('update:id', $event.target.value)"
		/>
	</FormRow>
	<div v-if="connection === 'tcpip'">
		<FormRow
			id="modbusHost"
			:label="$t('config.modbus.host')"
			:help="$t('config.modbus.hostHint')"
		>
			<PropertyField
				id="modbusHost"
				property="host"
				type="String"
				class="me-2"
				required
				:model-value="host"
				@change="$emit('update:host', $event.target.value)"
			/>
		</FormRow>
		<FormRow id="modbusPort" :label="$t('config.modbus.port')">
			<PropertyField
				id="modbusPort"
				property="port"
				type="Int"
				class="me-2 w-50"
				required
				:model-value="port || defaultPort || 502"
				@change="$emit('update:port', $event.target.value)"
			/>
		</FormRow>
		<FormRow
			v-if="showProtocolOptions"
			id="modbusTcp"
			:label="$t('config.modbus.protocol')"
			:help="
				protocol === 'tcp'
					? $t('config.modbus.protocolHintTcp')
					: $t('config.modbus.protocolHintRtu')
			"
		>
			<div class="btn-group" role="group">
				<input
					id="modbusTcp"
					v-model="protocol"
					type="radio"
					class="btn-check"
					name="modbusProtocol"
					value="tcp"
					tabindex="0"
					autocomplete="off"
				/>
				<label class="btn btn-outline-primary" for="modbusTcp">
					{{ $t("config.modbus.protocolValueTcp") }}
				</label>
				<input
					id="modbusRtu"
					v-model="protocol"
					type="radio"
					class="btn-check"
					name="modbusProtocol"
					value="rtu"
					tabindex="0"
					autocomplete="off"
				/>
				<label class="btn btn-outline-primary" for="modbusRtu">
					{{ $t("config.modbus.protocolValueRtu") }}
				</label>
			</div>
		</FormRow>
	</div>
	<div v-else>
		<FormRow
			id="modbusDevice"
			:label="$t('config.modbus.device')"
			:help="$t('config.modbus.deviceHint')"
		>
			<PropertyField
				id="modbusDevice"
				property="device"
				type="String"
				class="me-2"
				required
				:model-value="device"
				@change="$emit('update:device', $event.target.value)"
			/>
		</FormRow>
		<FormRow id="modbusBaudrate" :label="$t('config.modbus.baudrate')">
			<PropertyField
				id="modbusBaudrate"
				property="baudrate"
				type="Choice"
				class="me-2 w-50"
				:choice="baudrateOptions"
				required
				:model-value="baudrate || defaultBaudrate"
				@change="$emit('update:baudrate', $event.target.value)"
			/>
		</FormRow>
		<FormRow id="modbusComset" :label="$t('config.modbus.comset')">
			<PropertyField
				id="modbusComset"
				property="comset"
				type="Choice"
				class="me-2 w-50"
				:choice="comsetOptions"
				required
				:model-value="comset || defaultComset || '8N1'"
				@change="$emit('update:comset', $event.target.value)"
			/>
		</FormRow>
	</div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import FormRow from "../FormRow.vue";
import PropertyField from "../PropertyField.vue";
import type { PropType } from "vue";
import type { ModbusCapability } from "./index";
type Modbus = "rs485serial" | "rs485tcpip" | "tcpip";
type ConnectionOption = "tcpip" | "serial";
type ProtocolOption = "tcp" | "rtu";

export default defineComponent({
	name: "Modbus",
	components: { FormRow, PropertyField },
	props: {
		capabilities: {
			type: Array as PropType<ModbusCapability[]>,
			default: () => [],
		},
		modbus: String as PropType<Modbus>,
		host: String,
		port: [Number, String],
		id: [Number, String],
		baudrate: [Number, String],
		comset: String,
		device: String,
		defaultPort: Number,
		defaultId: Number,
		defaultComset: String,
		defaultBaudrate: Number,
	},
	emits: [
		"update:modbus",
		"update:id",
		"update:host",
		"update:port",
		"update:device",
		"update:baudrate",
		"update:comset",
	],
	data(): { connection: ConnectionOption; protocol: ProtocolOption } {
		return {
			connection: "tcpip",
			protocol: "tcp",
		};
	},
	computed: {
		selectedModbus(): Modbus {
			if (this.connection === "serial") {
				return "rs485serial";
			}
			return this.protocol === "rtu" ? "rs485tcpip" : "tcpip";
		},
		showConnectionOptions() {
			return this.capabilities.includes("rs485");
		},
		showProtocolOptions() {
			return this.connection === "tcpip" && this.capabilities.includes("rs485");
		},
		comsetOptions() {
			return ["8N1", "8E1", "8N2"].map((v) => {
				return { key: v, name: v };
			});
		},
		baudrateOptions() {
			return [1200, 9600, 19200, 38400, 57600, 115200].map((v) => {
				return { key: v, name: `${v}` };
			});
		},
	},
	watch: {
		selectedModbus(newValue: Modbus) {
			this.$emit("update:modbus", newValue);
		},
		options(newValue: ModbusCapability[]) {
			this.setProtocolByCapabilities(newValue);
			this.$emit("update:modbus", this.selectedModbus);
		},
		modbus(newValue: Modbus) {
			if (newValue) {
				this.setConnectionAndProtocolByModbus(newValue);
			}
		},
	},
	mounted() {
		this.setConnectionAndProtocolByModbus(this.modbus);
		this.$emit("update:modbus", this.selectedModbus);
	},
	methods: {
		setProtocolByCapabilities(capabilities: ModbusCapability[]) {
			this.protocol = capabilities.includes("tcpip") ? "tcp" : "rtu";
		},
		setConnectionAndProtocolByModbus(modbus?: Modbus) {
			switch (modbus) {
				case "rs485serial":
					this.connection = "serial";
					this.protocol = "rtu";
					break;
				case "rs485tcpip":
					this.connection = "tcpip";
					this.protocol = "rtu";
					break;
				case "tcpip":
					this.connection = "tcpip";
					this.protocol = "tcp";
					break;
			}
		},
	},
});
</script>
