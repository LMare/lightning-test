

/**
 * Add RGAA attributes on the lines of the tables
 */
function enrichRgaaTableRows(tbody, columnIds = []) {
	const rows = tbody.querySelectorAll("tr");
	let prefixId = tbody.id ? tbody.id : "row";
	rows.forEach((row, rowIndex) => {
		const th = row.querySelector("th");
		let rowId = "";
		if (th) {
			if(tbody.id){
				prefixId = `${prefixId}-row`;
			} else {
				console.warn("No id on the tbody", tbody);
			}
			rowId = `${prefixId}-${rowIndex}`;
			th.id = rowId;
		}

		const tds = row.querySelectorAll("td");
		tds.forEach((td, i) => {
			const colId = columnIds[i];
			if (colId) {
				td.setAttribute("headers", `${rowId} ${colId}`);
			}
		});
	});
}
