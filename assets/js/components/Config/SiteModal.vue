<template>
	<GenericModal
		id="siteModal"
		ref="modal"
		:title="$t('config.site.title')"
		data-testid="site-modal"
		config-modal-name="site"
		@open="open"
	>
		<p v-if="error" class="text-danger">{{ error }}</p>
		<form ref="form" class="container mx-0 px-0" @submit.prevent="save">
			<FormRow
				id="siteTitle"
				:label="$t('config.site.sitetitle.label')"
				:help="$t('config.site.sitetitle.description')"
			>
				<input id="siteTitle" v-model="title" class="form-control" />
			</FormRow>

			<FormRow 
				id="currency" 
				:label="$t('config.site.currency.label')"
				:example="exampleText"
				:help="$t('config.site.currency.description')"
			>
				<select id="currency" v-model="selectedCurrency" class="form-select" required>
					<option
						v-for="currency in currencies"
						:key="currency.code"
						:value="currency.code"
					>
						{{ currency.code }} - {{ currency.name }}
					</option>
				</select>
			</FormRow>

			<div class="mt-4 d-flex justify-content-between gap-2 flex-column flex-sm-row">
				<button
					type="button"
					class="btn btn-link text-muted btn-cancel"
					data-bs-dismiss="modal"
				>
					{{ $t('config.general.cancel') }}
				</button>

				<button
					type="submit"
					class="btn btn-primary order-1 order-sm-2 flex-grow-1 flex-sm-grow-0 px-4"
					:disabled="saving || !changed"
				>
					<span
						v-if="saving"
						class="spinner-border spinner-border-sm"
						role="status"
						aria-hidden="true"
					></span>
					{{ $t('config.general.save') }}
				</button>
			</div>
		</form>
	</GenericModal>
</template>

<script lang="ts">
import GenericModal from "../Helper/GenericModal.vue";
import FormRow from "./FormRow.vue";
import store from "@/store";
import api from "@/api";
import { CURRENCY } from "@/types/evcc";
import formatter from "@/mixins/formatter";

export default {
	name: "SiteModal",
	components: { FormRow, GenericModal },
	mixins: [formatter],
	emits: ["changed"],
	data() {
		return {
			saving: false,
			error: "",
			selectedCurrency: "EUR",
			initialCurrency: "EUR",
			title: "",
			initialTitle: "",
		};
	},
	computed: {
		currencies() {
			return Object.values(CURRENCY).map((code) => ({
				code,
				name: this.fmtCurrencyName(code),
			}));
		},
		changed() {
			return (
				this.title !== this.initialTitle ||
				this.selectedCurrency !== this.initialCurrency
			);
		},
		exampleText() {
			const price = this.fmtPricePerKWh(0.122, this.selectedCurrency);
			const amount = this.fmtMoney(20.2, this.selectedCurrency, true, true);
			return this.$t("config.site.currency.example", { price, amount });
		},
	},
	methods: {
		reset() {
			const currency = store?.state?.currency || "EUR";
			this.saving = false;
			this.error = "";
			this.selectedCurrency = currency;
			this.initialCurrency = currency;
			this.title = store.state?.siteTitle || "";
			this.initialTitle = this.title;
		},
		async open() {
			this.reset();
		},
		async save() {
			this.saving = true;
			this.error = "";
			try {
				const requests = [];
				if (this.title !== this.initialTitle) {
					requests.push(api.put("/config/site", { title: this.title }));
				}
				if (this.selectedCurrency !== this.initialCurrency) {
					requests.push(api.put("/config/currency", JSON.stringify(this.selectedCurrency)));
				}
				await Promise.all(requests);
				this.$emit("changed");
				this.$refs.modal.close();
			} catch (e) {
				this.error = (e && (e as any).message) || String(e);
			}
			this.saving = false;
		},
	},
};
</script>

<style scoped>
.container {
	margin-left: calc(var(--bs-gutter-x) * -0.5);
	margin-right: calc(var(--bs-gutter-x) * -0.5);
	padding-right: 0;
}
</style>
